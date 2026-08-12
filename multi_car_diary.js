// multi_car_diary.js — JavaScript версия

const fs = require('fs');
const readline = require('readline');

const DATA_FILE = 'cars.json';
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let cars = [];

function load() {
    if (fs.existsSync(DATA_FILE)) {
        try {
            cars = JSON.parse(fs.readFileSync(DATA_FILE, 'utf8'));
        } catch (e) {
            cars = [];
        }
    }
}

function save() {
    fs.writeFileSync(DATA_FILE, JSON.stringify(cars, null, 2));
}

function ask(question) {
    return new Promise(resolve => rl.question(question, resolve));
}

function color(text, code) {
    return `\x1b[${code}m${text}\x1b[0m`;
}

function findCar(id) {
    return cars.find(c => c.id === id);
}

function listCars() {
    if (cars.length === 0) {
        console.log(color("Нет автомобилей.", "33"));
        return;
    }
    console.log(color(`${'ID'.padEnd(4)} ${'Марка'.padEnd(15)} ${'Модель'.padEnd(15)} ${'Год'.padEnd(6)} ${'Пробег'.padEnd(10)}`, "36"));
    console.log("-".repeat(55));
    for (const c of cars) {
        console.log(`${String(c.id).padEnd(4)} ${c.brand.padEnd(15)} ${c.model.padEnd(15)} ${String(c.year).padEnd(6)} ${String(c.mileage).padEnd(10)}`);
    }
}

async function addRepair(car) {
    const date = await ask("Дата (ГГГГ-ММ-ДД): ");
    const mileage = parseInt(await ask("Пробег (км): "));
    const desc = await ask("Описание: ");
    const cost = parseFloat(await ask("Стоимость (руб): "));
    car.repairs.push({ date, mileage, description: desc, cost });
    if (mileage > car.mileage) car.mileage = mileage;
    save();
    console.log(color("✅ Ремонт добавлен.", "32"));
}

async function addOilChange(car) {
    const date = await ask("Дата (ГГГГ-ММ-ДД): ");
    const mileage = parseInt(await ask("Пробег (км): "));
    const oil = await ask("Тип масла: ");
    const filter = await ask("Фильтр (артикул): ");
    const cost = parseFloat(await ask("Стоимость (руб): "));
    const interval = parseInt(await ask("Интервал (км, по умолч. 10000): ")) || 10000;
    car.oil_changes.push({
        date, mileage, oil_type: oil, filter, cost, interval, next_mileage: mileage + interval
    });
    if (mileage > car.mileage) car.mileage = mileage;
    save();
    console.log(color("✅ Замена масла добавлена.", "32"));
}

async function addReminder(car) {
    const title = await ask("Название: ");
    const date = await ask("Дата (ГГГГ-ММ-ДД): ");
    const mileage = parseInt(await ask("Пробег (км): "));
    const desc = await ask("Описание: ");
    let priority = (await ask("Приоритет (низкий/средний/высокий): ")).toLowerCase();
    if (!["низкий", "средний", "высокий"].includes(priority)) priority = "средний";
    car.reminders.push({ title, date, mileage, description: desc, priority, done: false });
    save();
    console.log(color("✅ Напоминание добавлено.", "32"));
}

function showHistory(car) {
    console.log(`\n${color(`📋 История для ${car.brand} ${car.model} ${car.year}`, "36")}`);
    if (car.repairs.length > 0) {
        console.log(color("Ремонты:", "33"));
        for (const r of car.repairs) {
            console.log(`  ${r.date} | ${r.mileage} км | ${r.description} | ${r.cost} руб.`);
        }
    }
    if (car.oil_changes.length > 0) {
        console.log(color("Замены масла:", "33"));
        for (const o of car.oil_changes) {
            console.log(`  ${o.date} | ${o.mileage} км | ${o.oil_type} | ${o.filter} | ${o.cost} руб. | след. ${o.next_mileage} км`);
        }
    }
    if (car.reminders.length > 0) {
        console.log(color("Напоминания:", "33"));
        for (const r of car.reminders) {
            const status = r.done ? "✅" : "⏳";
            console.log(`  ${status} ${r.title} | ${r.date} | ${r.priority} | ${r.description}`);
        }
    }
}

function showStats(car) {
    const totalRepair = car.repairs.reduce((s, r) => s + r.cost, 0);
    const totalOil = car.oil_changes.reduce((s, o) => s + o.cost, 0);
    const total = totalRepair + totalOil;
    console.log(`\n${color(`📊 Статистика для ${car.brand} ${car.model} ${car.year}`, "36")}`);
    console.log(`  Всего ремонтов: ${car.repairs.length}`);
    console.log(`  Всего замен масла: ${car.oil_changes.length}`);
    console.log(`  Всего напоминаний: ${car.reminders.length}`);
    console.log(`  Общая стоимость: ${total.toFixed(2)} руб.`);
    if (car.repairs.length > 0) {
        console.log(`  Средняя стоимость ремонта: ${(totalRepair / car.repairs.length).toFixed(2)} руб.`);
    }
}

async function carMenu(car) {
    while (true) {
        console.log(`\n${color(`🚗 ${car.brand} ${car.model} ${car.year} (Пробег: ${car.mileage} км) — меню автомобиля`, "36")}`);
        console.log("1. Добавить ремонт");
        console.log("2. Добавить замену масла");
        console.log("3. Добавить напоминание");
        console.log("4. Показать историю");
        console.log("5. Показать статистику");
        console.log("6. Назад");
        const choice = await ask("Выберите действие: ");
        switch (choice.trim()) {
            case "1": await addRepair(car); break;
            case "2": await addOilChange(car); break;
            case "3": await addReminder(car); break;
            case "4": showHistory(car); break;
            case "5": showStats(car); break;
            case "6": return;
            default: console.log(color("Неверный выбор.", "31"));
        }
    }
}

async function main() {
    load();
    while (true) {
        console.log(`\n${color("🚗 Многомарочный дневник автомобиля (JavaScript)", "36")}`);
        console.log("1. Выбрать/создать автомобиль");
        console.log("2. Удалить автомобиль");
        console.log("3. Показать все автомобили");
        console.log("4. Выход");
        const choice = await ask("Выберите действие: ");
        switch (choice.trim()) {
            case "1": {
                listCars();
                if (cars.length > 0) {
                    const idStr = await ask("Введите номер автомобиля (или 0 для создания нового): ");
                    if (idStr.trim() === "0") {
                        const brand = await ask("Марка: ");
                        const model = await ask("Модель: ");
                        const year = parseInt(await ask("Год: "));
                        const vin = await ask("VIN: ");
                        const mileage = parseInt(await ask("Текущий пробег (км): "));
                        const id = cars.length + 1;
                        cars.push({ id, brand, model, year, vin, mileage, repairs: [], oil_changes: [], reminders: [] });
                        save();
                        console.log(color(`✅ Автомобиль добавлен (ID: ${id})`, "32"));
                        const car = findCar(id);
                        if (car) await carMenu(car);
                    } else {
                        const id = parseInt(idStr);
                        const car = findCar(id);
                        if (car) {
                            await carMenu(car);
                        } else {
                            console.log(color("❌ Автомобиль не найден.", "31"));
                        }
                    }
                } else {
                    const brand = await ask("Марка: ");
                    const model = await ask("Модель: ");
                    const year = parseInt(await ask("Год: "));
                    const vin = await ask("VIN: ");
                    const mileage = parseInt(await ask("Текущий пробег (км): "));
                    const id = cars.length + 1;
                    cars.push({ id, brand, model, year, vin, mileage, repairs: [], oil_changes: [], reminders: [] });
                    save();
                    console.log(color(`✅ Автомобиль добавлен (ID: ${id})`, "32"));
                    const car = findCar(id);
                    if (car) await carMenu(car);
                }
                break;
            }
            case "2": {
                listCars();
                const id = parseInt(await ask("Введите ID автомобиля для удаления: "));
                const index = cars.findIndex(c => c.id === id);
                if (index !== -1) {
                    cars.splice(index, 1);
                    save();
                    console.log(color("✅ Автомобиль удалён.", "32"));
                } else {
                    console.log(color("❌ Автомобиль не найден.", "31"));
                }
                break;
            }
            case "3": listCars(); break;
            case "4": console.log("До свидания!"); rl.close(); return;
            default: console.log(color("Неверный выбор.", "31"));
        }
    }
}

main().catch(console.error);
