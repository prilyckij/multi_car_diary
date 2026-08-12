// multi_car_diary.rs — Rust версия

use serde::{Deserialize, Serialize};
use std::fs;
use std::io::{self, Write};

#[derive(Serialize, Deserialize, Clone)]
struct Repair {
    date: String,
    mileage: u32,
    description: String,
    cost: f64,
}

#[derive(Serialize, Deserialize, Clone)]
struct OilChange {
    date: String,
    mileage: u32,
    oil_type: String,
    filter: String,
    cost: f64,
    interval: u32,
    next_mileage: u32,
}

#[derive(Serialize, Deserialize, Clone)]
struct Reminder {
    title: String,
    date: String,
    mileage: u32,
    description: String,
    priority: String,
    done: bool,
}

#[derive(Serialize, Deserialize, Clone)]
struct Car {
    id: usize,
    brand: String,
    model: String,
    year: u32,
    vin: String,
    mileage: u32,
    repairs: Vec<Repair>,
    oil_changes: Vec<OilChange>,
    reminders: Vec<Reminder>,
}

struct Manager {
    cars: Vec<Car>,
    file: String,
}

impl Manager {
    fn new(file: &str) -> Self {
        let mut m = Manager { cars: Vec::new(), file: file.to_string() };
        m.load();
        m
    }

    fn load(&mut self) {
        if let Ok(data) = fs::read_to_string(&self.file) {
            if let Ok(cars) = serde_json::from_str(&data) {
                self.cars = cars;
                return;
            }
        }
        self.cars = Vec::new();
    }

    fn save(&self) {
        let data = serde_json::to_string_pretty(&self.cars).unwrap();
        fs::write(&self.file, data).unwrap();
    }

    fn add_car(&mut self, brand: String, model: String, year: u32, vin: String, mileage: u32) -> usize {
        let id = self.cars.len() + 1;
        self.cars.push(Car {
            id,
            brand,
            model,
            year,
            vin,
            mileage,
            repairs: Vec::new(),
            oil_changes: Vec::new(),
            reminders: Vec::new(),
        });
        self.save();
        id
    }

    fn delete_car(&mut self, id: usize) -> bool {
        let pos = self.cars.iter().position(|c| c.id == id);
        if let Some(idx) = pos {
            self.cars.remove(idx);
            self.save();
            true
        } else {
            false
        }
    }

    fn get_car(&self, id: usize) -> Option<&Car> {
        self.cars.iter().find(|c| c.id == id)
    }

    fn get_car_mut(&mut self, id: usize) -> Option<&mut Car> {
        self.cars.iter_mut().find(|c| c.id == id)
    }

    fn list_cars(&self) {
        if self.cars.is_empty() {
            println!("\x1b[33mНет автомобилей.\x1b[0m");
            return;
        }
        println!("\x1b[36m{:<4} {:<15} {:<15} {:<6} {:<10}\x1b[0m", "ID", "Марка", "Модель", "Год", "Пробег");
        println!("{}", "-".repeat(55));
        for c in &self.cars {
            println!("{:<4} {:<15} {:<15} {:<6} {:<10}", c.id, c.brand, c.model, c.year, c.mileage);
        }
    }
}

fn main() {
    let mut manager = Manager::new("cars.json");
    loop {
        println!("\n\x1b[36m🚗 Многомарочный дневник автомобиля (Rust)\x1b[0m");
        println!("1. Выбрать/создать автомобиль");
        println!("2. Удалить автомобиль");
        println!("3. Показать все автомобили");
        println!("4. Выход");
        print!("Выберите действие: ");
        io::stdout().flush().unwrap();
        let mut choice = String::new();
        io::stdin().read_line(&mut choice).unwrap();
        match choice.trim() {
            "1" => {
                manager.list_cars();
                if !manager.cars.is_empty() {
                    print!("Введите номер автомобиля (или 0 для создания нового): ");
                    io::stdout().flush().unwrap();
                    let mut input = String::new();
                    io::stdin().read_line(&mut input).unwrap();
                    let input = input.trim();
                    if input == "0" {
                        create_new_car(&mut manager);
                    } else if let Ok(id) = input.parse::<usize>() {
                        if let Some(car) = manager.get_car(id) {
                            car_menu(&mut manager, car.id);
                        } else {
                            println!("\x1b[31m❌ Автомобиль не найден.\x1b[0m");
                        }
                    }
                } else {
                    create_new_car(&mut manager);
                }
            }
            "2" => {
                manager.list_cars();
                print!("Введите ID автомобиля для удаления: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                if let Ok(id) = id_str.trim().parse::<usize>() {
                    if manager.delete_car(id) {
                        println!("\x1b[32m✅ Автомобиль удалён.\x1b[0m");
                    } else {
                        println!("\x1b[31m❌ Автомобиль не найден.\x1b[0m");
                    }
                }
            }
            "3" => manager.list_cars(),
            "4" => {
                println!("До свидания!");
                break;
            }
            _ => println!("\x1b[31mНеверный выбор.\x1b[0m"),
        }
    }
}

fn create_new_car(manager: &mut Manager) {
    print!("Марка: ");
    io::stdout().flush().unwrap();
    let mut brand = String::new();
    io::stdin().read_line(&mut brand).unwrap();
    let brand = brand.trim().to_string();
    print!("Модель: ");
    io::stdout().flush().unwrap();
    let mut model = String::new();
    io::stdin().read_line(&mut model).unwrap();
    let model = model.trim().to_string();
    print!("Год: ");
    io::stdout().flush().unwrap();
    let mut year_str = String::new();
    io::stdin().read_line(&mut year_str).unwrap();
    let year: u32 = year_str.trim().parse().unwrap();
    print!("VIN: ");
    io::stdout().flush().unwrap();
    let mut vin = String::new();
    io::stdin().read_line(&mut vin).unwrap();
    let vin = vin.trim().to_string();
    print!("Текущий пробег (км): ");
    io::stdout().flush().unwrap();
    let mut mileage_str = String::new();
    io::stdin().read_line(&mut mileage_str).unwrap();
    let mileage: u32 = mileage_str.trim().parse().unwrap();
    let id = manager.add_car(brand, model, year, vin, mileage);
    println!("\x1b[32m✅ Автомобиль добавлен (ID: {})\x1b[0m", id);
    car_menu(manager, id);
}

fn car_menu(manager: &mut Manager, car_id: usize) {
    loop {
        let car = manager.get_car(car_id).unwrap();
        println!("\n\x1b[36m🚗 {} {} {} (Пробег: {} км) — меню автомобиля\x1b[0m", car.brand, car.model, car.year, car.mileage);
        println!("1. Добавить ремонт");
        println!("2. Добавить замену масла");
        println!("3. Добавить напоминание");
        println!("4. Показать историю");
        println!("5. Показать статистику");
        println!("6. Назад");
        print!("Выберите действие: ");
        io::stdout().flush().unwrap();
        let mut choice = String::new();
        io::stdin().read_line(&mut choice).unwrap();
        match choice.trim() {
            "1" => add_repair(manager, car_id),
            "2" => add_oil_change(manager, car_id),
            "3" => add_reminder(manager, car_id),
            "4" => show_history(manager, car_id),
            "5" => show_stats(manager, car_id),
            "6" => break,
            _ => println!("\x1b[31mНеверный выбор.\x1b[0m"),
        }
    }
}

fn add_repair(manager: &mut Manager, car_id: usize) {
    print!("Дата (ГГГГ-ММ-ДД): ");
    io::stdout().flush().unwrap();
    let mut date = String::new();
    io::stdin().read_line(&mut date).unwrap();
    let date = date.trim().to_string();
    print!("Пробег (км): ");
    io::stdout().flush().unwrap();
    let mut mileage_str = String::new();
    io::stdin().read_line(&mut mileage_str).unwrap();
    let mileage: u32 = mileage_str.trim().parse().unwrap();
    print!("Описание: ");
    io::stdout().flush().unwrap();
    let mut desc = String::new();
    io::stdin().read_line(&mut desc).unwrap();
    let desc = desc.trim().to_string();
    print!("Стоимость (руб): ");
    io::stdout().flush().unwrap();
    let mut cost_str = String::new();
    io::stdin().read_line(&mut cost_str).unwrap();
    let cost: f64 = cost_str.trim().parse().unwrap();
    let car = manager.get_car_mut(car_id).unwrap();
    car.repairs.push(Repair { date, mileage, description: desc, cost });
    if mileage > car.mileage { car.mileage = mileage; }
    manager.save();
    println!("\x1b[32m✅ Ремонт добавлен.\x1b[0m");
}

fn add_oil_change(manager: &mut Manager, car_id: usize) {
    print!("Дата (ГГГГ-ММ-ДД): ");
    io::stdout().flush().unwrap();
    let mut date = String::new();
    io::stdin().read_line(&mut date).unwrap();
    let date = date.trim().to_string();
    print!("Пробег (км): ");
    io::stdout().flush().unwrap();
    let mut mileage_str = String::new();
    io::stdin().read_line(&mut mileage_str).unwrap();
    let mileage: u32 = mileage_str.trim().parse().unwrap();
    print!("Тип масла: ");
    io::stdout().flush().unwrap();
    let mut oil = String::new();
    io::stdin().read_line(&mut oil).unwrap();
    let oil = oil.trim().to_string();
    print!("Фильтр (артикул): ");
    io::stdout().flush().unwrap();
    let mut filter = String::new();
    io::stdin().read_line(&mut filter).unwrap();
    let filter = filter.trim().to_string();
    print!("Стоимость (руб): ");
    io::stdout().flush().unwrap();
    let mut cost_str = String::new();
    io::stdin().read_line(&mut cost_str).unwrap();
    let cost: f64 = cost_str.trim().parse().unwrap();
    print!("Интервал (км, по умолч. 10000): ");
    io::stdout().flush().unwrap();
    let mut interval_str = String::new();
    io::stdin().read_line(&mut interval_str).unwrap();
    let interval = if interval_str.trim().is_empty() { 10000 } else { interval_str.trim().parse().unwrap() };
    let next_mileage = mileage + interval;
    let car = manager.get_car_mut(car_id).unwrap();
    car.oil_changes.push(OilChange { date, mileage, oil_type: oil, filter, cost, interval, next_mileage });
    if mileage > car.mileage { car.mileage = mileage; }
    manager.save();
    println!("\x1b[32m✅ Замена масла добавлена.\x1b[0m");
}

fn add_reminder(manager: &mut Manager, car_id: usize) {
    print!("Название: ");
    io::stdout().flush().unwrap();
    let mut title = String::new();
    io::stdin().read_line(&mut title).unwrap();
    let title = title.trim().to_string();
    print!("Дата (ГГГГ-ММ-ДД): ");
    io::stdout().flush().unwrap();
    let mut date = String::new();
    io::stdin().read_line(&mut date).unwrap();
    let date = date.trim().to_string();
    print!("Пробег (км): ");
    io::stdout().flush().unwrap();
    let mut mileage_str = String::new();
    io::stdin().read_line(&mut mileage_str).unwrap();
    let mileage: u32 = mileage_str.trim().parse().unwrap();
    print!("Описание: ");
    io::stdout().flush().unwrap();
    let mut desc = String::new();
    io::stdin().read_line(&mut desc).unwrap();
    let desc = desc.trim().to_string();
    print!("Приоритет (низкий/средний/высокий): ");
    io::stdout().flush().unwrap();
    let mut priority = String::new();
    io::stdin().read_line(&mut priority).unwrap();
    let priority = priority.trim().to_lowercase();
    let priority = if priority == "низкий" || priority == "средний" || priority == "высокий" { priority } else { "средний".to_string() };
    let car = manager.get_car_mut(car_id).unwrap();
    car.reminders.push(Reminder { title, date, mileage, description: desc, priority, done: false });
    manager.save();
    println!("\x1b[32m✅ Напоминание добавлено.\x1b[0m");
}

fn show_history(manager: &Manager, car_id: usize) {
    let car = manager.get_car(car_id).unwrap();
    println!("\n\x1b[36m📋 История для {} {} {}\x1b[0m", car.brand, car.model, car.year);
    if !car.repairs.is_empty() {
        println!("\x1b[33mРемонты:\x1b[0m");
        for r in &car.repairs {
            println!("  {} | {} км | {} | {} руб.", r.date, r.mileage, r.description, r.cost);
        }
    }
    if !car.oil_changes.is_empty() {
        println!("\x1b[33mЗамены масла:\x1b[0m");
        for o in &car.oil_changes {
            println!("  {} | {} км | {} | {} | {} руб. | след. {} км", o.date, o.mileage, o.oil_type, o.filter, o.cost, o.next_mileage);
        }
    }
    if !car.reminders.is_empty() {
        println!("\x1b[33mНапоминания:\x1b[0m");
        for r in &car.reminders {
            let status = if r.done { "✅" } else { "⏳" };
            println!("  {} {} | {} | {} | {}", status, r.title, r.date, r.priority, r.description);
        }
    }
}

fn show_stats(manager: &Manager, car_id: usize) {
    let car = manager.get_car(car_id).unwrap();
    let total_repair: f64 = car.repairs.iter().map(|r| r.cost).sum();
    let total_oil: f64 = car.oil_changes.iter().map(|o| o.cost).sum();
    let total = total_repair + total_oil;
    println!("\n\x1b[36m📊 Статистика для {} {} {}\x1b[0m", car.brand, car.model, car.year);
    println!("  Всего ремонтов: {}", car.repairs.len());
    println!("  Всего замен масла: {}", car.oil_changes.len());
    println!("  Всего напоминаний: {}", car.reminders.len());
    println!("  Общая стоимость: {:.2} руб.", total);
    if !car.repairs.is_empty() {
        println!("  Средняя стоимость ремонта: {:.2} руб.", total_repair / car.repairs.len() as f64);
    }
}
