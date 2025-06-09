document.addEventListener("DOMContentLoaded", function () {
  const eventsContainer = document.getElementById("events");
  let eventSource = null;
  let reconnectAttempts = 0;
  const maxReconnectAttempts = 5;
  const reconnectDelay = 3000; // 3 segundos

  function connectSSE() {
    if (eventSource) {
      eventSource.close();
    }

    console.log("Conectando a SSE...");
    eventSource = new EventSource("/events");

    eventSource.onopen = function () {
      console.log("Conexión SSE establecida");
      reconnectAttempts = 0;
    };

    eventSource.onmessage = function (event) {
      console.log("Evento recibido:", event.data);

      if (event.data === "ping") {
        console.log("Ping recibido");
        return;
      }

      if (event.data === "connected") {
        console.log("Conexión SSE confirmada");
        return;
      }

      try {
        const eventData = JSON.parse(event.data);
        const eventElement = document.createElement("div");
        eventElement.className = "border rounded p-4 bg-gray-50";

        eventElement.innerHTML = `
                <div class="flex justify-between items-center mb-2">
                    <div class="flex items-center space-x-2">
                        <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
                            ${eventData.type}
                        </span>
                        <span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">
                            ${eventData.status}
                        </span>
                    </div>
                    <span class="text-sm text-gray-500">${new Date(
                      eventData.timestamp
                    ).toLocaleTimeString()}</span>
                </div>
                <div class="mt-2">
                    <div class="text-xs text-gray-500 mb-1">ID: ${
                      eventData.id
                    }</div>
                    <pre class="text-sm bg-white p-2 rounded border">${
                      eventData.payload
                    }</pre>
                </div>
            `;

        eventsContainer.insertBefore(eventElement, eventsContainer.firstChild);
      } catch (error) {
        console.error("Error al procesar el evento:", error);
      }
    };

    eventSource.onerror = function (error) {
      console.error("Error en la conexión SSE:", error);

      if (eventSource) {
        eventSource.close();
        eventSource = null;
      }

      if (reconnectAttempts < maxReconnectAttempts) {
        reconnectAttempts++;
        console.log(
          `Intentando reconectar (${reconnectAttempts}/${maxReconnectAttempts})...`
        );
        setTimeout(connectSSE, reconnectDelay);
      } else {
        console.error("Número máximo de intentos de reconexión alcanzado");
        // Mostrar un mensaje al usuario
        const errorElement = document.createElement("div");
        errorElement.className =
          "bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded";
        errorElement.textContent =
          "Error de conexión. Por favor, recarga la página.";
        eventsContainer.insertBefore(errorElement, eventsContainer.firstChild);
      }
    };
  }

  // Iniciar la conexión
  connectSSE();

  // Limpiar la conexión cuando se cierre la página
  window.addEventListener("beforeunload", function () {
    if (eventSource) {
      eventSource.close();
    }
  });
});
