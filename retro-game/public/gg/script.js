function initGameAfterAuth(event) {
    if (event) event.preventDefault(); // Останавливаем перезагрузку

    // 1. Собираем данные из полей (ID должны совпадать с HTML)
    const fName = document.getElementById("firstName").value.trim();
    const lName = document.getElementById("lastName").value.trim();
    const phone = document.getElementById("phone").value.trim();

    if (!fName || !phone) {
        alert("⚠️ Введите имя и телефон!");
        return;
    }

    // 2. Отправляем данные в твой Go-сервер (в MySQL/SQLite)
    const formData = new FormData();
    formData.append('firstName', fName);
    formData.append('lastName', lName);
    formData.append('phone', phone);

    fetch('/submit_data', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            console.log("Данные успешно записаны в базу");
            
            // 3. ЗАПУСК ИГРЫ (Оригинальная функция из твоего HTML)
            // Мы вызываем ту самую функцию, которая прописана в теге <script> внутри HTML
            // Она переключит экраны и установит game.isProcessing = false
            launchOriginalGame(fName, lName, phone);
        } else {
            alert("Ошибка сохранения в БД. Проверь Go-сервер.");
        }
    })
    .catch(err => {
        console.error("Ошибка связи с сервером:", err);
        // Если сервер Go не запущен, всё равно даем поиграть (для тестов)
        launchOriginalGame(fName, lName, phone);
    });
}

function launchOriginalGame(fName, lName, phone) {
    // Скрываем панель входа
    document.getElementById("authPanel").style.display = "none";
    document.getElementById("gameScreen").style.display = "grid";

    // Устанавливаем имя игрока в интерфейсе
    const fullName = lName ? `${fName} ${lName}` : fName;
    document.getElementById("playerNameLabel").textContent = fullName;

    // СБРОС БЛОКИРОВКИ КАРТ
    // В твоем HTML объекте 'game' нужно выставить начальные значения
    if (typeof game !== 'undefined') {
        game.active = true;
        game.isProcessing = false; // Это самое важное, чтобы карты нажимались!
        game.user = { name: fullName };
        
        // Перерисовываем UI (функции из твоего HTML)
        if (typeof updateUI === 'function') updateUI();
        if (typeof addLog === 'function') addLog(`🔐 Доступ разрешен. Удачи, ${fName}!`, "#aaf0ff");
    }
}