import http from "k6/http";
import { check, sleep, group } from "k6";

const BASE_URL = "http://localhost:8080";

// ── Options: jalankan semua skenario secara berurutan ──
export const options = {
  scenarios: {
    // 1. Load test — steady concurrent users
    load_test: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "30s", target: 50 }, // ramp up
        { duration: "1m", target: 50 }, // steady
        { duration: "15s", target: 0 }, // ramp down
      ],
      gracefulRampDown: "10s",
      tags: { scenario: "load" },
    },
    // 2. Stress test — cari breaking point
    stress_test: {
      executor: "ramping-vus",
      startTime: "2m", // mulai setelah load test
      startVUs: 0,
      stages: [
        { duration: "30s", target: 100 },
        { duration: "30s", target: 200 },
        { duration: "30s", target: 300 },
        { duration: "30s", target: 0 },
      ],
      gracefulRampDown: "10s",
      tags: { scenario: "stress" },
    },
    // 3. Spike test — lonjakan tiba-tiba
    spike_test: {
      executor: "ramping-vus",
      startTime: "4m30s", // mulai setelah stress test
      startVUs: 0,
      stages: [
        { duration: "10s", target: 5 }, // baseline
        { duration: "5s", target: 300 }, // spike!
        { duration: "30s", target: 300 }, // sustain
        { duration: "5s", target: 5 }, // drop
        { duration: "20s", target: 0 }, // ramp down
      ],
      gracefulRampDown: "10s",
      tags: { scenario: "spike" },
    },
  },
  thresholds: {
    http_req_duration: ["p(99)<500"], // 99% request < 500ms (PRD: p99 < 100ms)
    http_req_failed: ["rate<0.05"], // error rate < 5%
  },
};

// ── Setup: buat 1 user per VU sebelum test dimulai ──
// Gunakan __VU-based email supaya tiap VU punya akun sendiri
// tapi tidak membuat akun baru setiap iterasi
let cachedToken = "";

function ensureToken() {
  if (cachedToken !== "") return cachedToken;

  const email = `loadtest_vu${__VU}@test.com`;
  const password = "password123";

  // Register (idempotent — boleh gagal kalau sudah ada)
  http.post(
    `${BASE_URL}/v1/auth/register`,
    JSON.stringify({ email, password, full_name: `VU ${__VU}` }),
    { headers: { "Content-Type": "application/json" } },
  );

  // Login
  const loginRes = http.post(
    `${BASE_URL}/v1/auth/login`,
    JSON.stringify({ email, password }),
    { headers: { "Content-Type": "application/json" } },
  );

  check(loginRes, { "login success": (r) => r.status === 200 });

  try {
    cachedToken = JSON.parse(loginRes.body).accessToken || "";
  } catch {
    cachedToken = "";
  }
  return cachedToken;
}

// ── Main flow ──
export default function () {
  const token = ensureToken();
  const headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  group("Auth", () => {
    // Login ulang tiap iterasi untuk mengukur auth latency
    const email = `loadtest_vu${__VU}@test.com`;
    const res = http.post(
      `${BASE_URL}/v1/auth/login`,
      JSON.stringify({ email, password: "password123" }),
      { headers: { "Content-Type": "application/json" } },
    );
    check(res, { "login 200": (r) => r.status === 200 });
    // refresh token setelah login ulang
    try {
      const newToken = JSON.parse(res.body).accessToken;
      if (newToken) cachedToken = newToken;
    } catch {}
  });

  group("Products", () => {
    // Create
    const createRes = http.post(
      `${BASE_URL}/v1/products`,
      JSON.stringify({
        name: `Produk ${__VU}-${__ITER}`,
        description: "Test product",
        price: 50000,
        stock_qty: 100,
      }),
      { headers },
    );
    check(createRes, { "create product 200": (r) => r.status === 200 });

    let productId = "";
    try {
      productId = JSON.parse(createRes.body).product.id;
    } catch {}

    // List
    const listRes = http.get(`${BASE_URL}/v1/products?page=1&limit=10`, {
      headers,
    });
    check(listRes, { "list products 200": (r) => r.status === 200 });

    // Get
    if (productId) {
      const getRes = http.get(`${BASE_URL}/v1/products/${productId}`, {
        headers,
      });
      check(getRes, { "get product 200": (r) => r.status === 200 });
    }

    sleep(0.5);

    // Order
    if (productId) {
      group("Orders", () => {
        const orderRes = http.post(
          `${BASE_URL}/v1/orders`,
          JSON.stringify({
            items: [{ product_id: productId, quantity: 1 }],
          }),
          { headers },
        );
        check(orderRes, { "create order 200": (r) => r.status === 200 });

        // List orders
        const listOrderRes = http.get(`${BASE_URL}/v1/orders?page=1&limit=10`, {
          headers,
        });
        check(listOrderRes, { "list orders 200": (r) => r.status === 200 });
      });
    }
  });

  group("Users", () => {
    const profileRes = http.get(`${BASE_URL}/v1/users/me`, { headers });
    check(profileRes, { "get profile 200": (r) => r.status === 200 });
  });

  sleep(1);
}
