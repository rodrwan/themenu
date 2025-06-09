document.addEventListener("DOMContentLoaded", () => {
  const updateOrderStatusForm = document.getElementById("updateOrderStatus");
  const orderIDInput = document.getElementById("orderID");
  const orderStatusSelect = document.getElementById("orderStatus");

  updateOrderStatusForm?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const orderID = orderIDInput.value;
    const status = orderStatusSelect.value;

    try {
      const response = await fetch(`/orders/${orderID}/status`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ status }),
      });

      if (!response.ok) {
        throw new Error("Error al actualizar el estado de la orden");
      }

      const data = await response.json();
      alert(data.message);
      orderIDInput.value = "";
      orderStatusSelect.value = "pending";
    } catch (error) {
      console.error("Error:", error);
      alert("Error al actualizar el estado de la orden");
    }
  });

  const eventSource = new EventSource("/events");
  const eventsContainer = document.getElementById("events");
  const ordersContainer = document.getElementById("orders");

  eventSource.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const eventElement = document.createElement("div");
    eventElement.className = "border rounded p-4 bg-gray-50";
    eventElement.innerHTML = `
      <div class="flex justify-between items-center mb-2">
        <div class="flex items-center space-x-2">
          <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
            ${data.type}
          </span>
          <span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">
            ${data.status}
          </span>
        </div>
        <span class="text-sm text-gray-500">${new Date(
          data.timestamp
        ).toLocaleTimeString()}</span>
      </div>
      <div class="mt-2">
        <div class="text-xs text-gray-500 mb-1">ID: ${data.id}</div>
        <pre class="text-sm bg-white p-2 rounded border">${data.payload}</pre>
      </div>
    `;
    eventsContainer.insertBefore(eventElement, eventsContainer.firstChild);
  };

  // Obtener órdenes activas desde el endpoint GET /orders
  fetch("/orders", {
    headers: {
      Authorization: "Bearer aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    },
  })
    .then((response) => response.json())
    .then((orders) => {
      orders.forEach((order) => {
        console.log(order);
        const orderElement = document.createElement("div");
        orderElement.className = "border rounded p-4 bg-gray-50";
        orderElement.innerHTML = `
        <div class="flex justify-between items-center mb-2">
          <div class="flex items-center space-x-2">
            <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
              Orden #${order.id}
            </span>
            <span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">
              ${order.status}
            </span>
          </div>
          <span class="text-sm text-gray-500">${new Date(
            order.created_at
          ).toLocaleString()}</span>
        </div>
        <div class="mt-2">
          <div class="text-xs text-gray-500 mb-1">ID: ${order.id}</div>
          <pre class="text-sm bg-white p-2 rounded border">${
            order.dish_name
          }</pre>
        </div>
        <div class="mt-4">
          <select class="orderStatus border rounded p-2" required>
            <option value="received">Recibido</option>
            <option value="confirmed">Confirmado</option>
            <option value="preparing">Preparando</option>
            <option value="served">Servido</option>
            <option value="cancelled">Cancelado</option>
          </select>
          <button type="button" class="updateOrderStatus bg-blue-500 text-white px-4 py-2 rounded ml-2">Actualizar Estado</button>
        </div>
      `;
        ordersContainer.appendChild(orderElement);
      });
    })
    .catch((error) => {
      console.error("Error:", error);
    });

  // Manejar la actualización del estado de las órdenes
  ordersContainer.addEventListener("click", function (event) {
    if (event.target.classList.contains("updateOrderStatus")) {
      const orderElement = event.target.closest(".border");
      const orderID = orderElement
        .querySelector(".text-xs.text-gray-500")
        .textContent.split(": ")[1];
      const orderStatus = orderElement.querySelector(".orderStatus").value;

      fetch(`/orders/${orderID}/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
        },
        body: JSON.stringify({ status: orderStatus }),
      })
        .then((response) => response.json())
        .then((data) => {
          alert(data.message);
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    }
  });
});
