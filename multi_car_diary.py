# multi_car_diary.py — Python версия

import json
import os
from datetime import datetime
from colorama import init, Fore, Style

init(autoreset=True)
DATA_FILE = "cars.json"

class Car:
    def __init__(self, id, brand, model, year, vin, mileage, repairs=None, oil_changes=None, reminders=None):
        self.id = id
        self.brand = brand
        self.model = model
        self.year = year
        self.vin = vin
        self.mileage = mileage
        self.repairs = repairs or []
        self.oil_changes = oil_changes or []
        self.reminders = reminders or []

    def to_dict(self):
        return {
            "id": self.id,
            "brand": self.brand,
            "model": self.model,
            "year": self.year,
            "vin": self.vin,
            "mileage": self.mileage,
            "repairs": self.repairs,
            "oil_changes": self.oil_changes,
            "reminders": self.reminders
        }

    @classmethod
    def from_dict(cls, data):
        return cls(
            data["id"], data["brand"], data["model"], data["year"],
            data["vin"], data["mileage"], data.get("repairs", []),
            data.get("oil_changes", []), data.get("reminders", [])
        )

    def add_repair(self, date, mileage, description, cost):
        self.repairs.append({
            "date": date,
            "mileage": mileage,
            "description": description,
            "cost": cost
        })
        self.mileage = max(self.mileage, mileage)

    def add_oil_change(self, date, mileage, oil_type, filter_ref, cost, interval=10000):
        next_mileage = mileage + interval
        self.oil_changes.append({
            "date": date,
            "mileage": mileage,
            "oil_type": oil_type,
            "filter": filter_ref,
            "cost": cost,
            "interval": interval,
            "next_mileage": next_mileage
        })
        self.mileage = max(self.mileage, mileage)

    def add_reminder(self, title, date, mileage, description, priority):
        self.reminders.append({
            "title": title,
            "date": date,
            "mileage": mileage,
            "description": description,
            "priority": priority,
            "done": False
        })

    def total_cost(self):
        return sum(r["cost"] for r in self.repairs) + sum(o["cost"] for o in self.oil_changes)

    def __str__(self):
        return f"{self.brand} {self.model} {self.year} (Пробег: {self.mileage} км)"

class CarManager:
    def __init__(self):
        self.cars = []
        self.current_car = None
        self.load()

    def load(self):
        if os.path.exists(DATA_FILE):
            try:
                with open(DATA_FILE, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    self.cars = [Car.from_dict(c) for c in data]
            except:
                self.cars = []

    def save(self):
        with open(DATA_FILE, 'w', encoding='utf-8') as f:
            json.dump([c.to_dict() for c in self.cars], f, indent=2, ensure_ascii=False)

    def add_car(self, brand, model, year, vin, mileage):
        id = len(self.cars) + 1
        car = Car(id, brand, model, year, vin, mileage)
        self.cars.append(car)
        self.save()
        return id

    def delete_car(self, id):
        for i, c in enumerate(self.cars):
            if c.id == id:
                del self.cars[i]
                self.save()
                return True
        return False

    def get_car(self, id):
        for c in self.cars:
            if c.id == id:
                return c
        return None

    def list_cars(self):
        if not self.cars:
            print(Fore.YELLOW + "Нет автомобилей.")
            return
        print(Fore.CYAN + f"{'ID':<4} {'Марка':<15} {'Модель':<15} {'Год':<6} {'Пробег':<10}")
        print("-" * 55)
        for c in self.cars:
            print(f"{c.id:<4} {c.brand:<15} {c.model:<15} {c.year:<6} {c.mileage:<10}")

    def car_menu(self, car):
        while True:
            print(f"\n{Fore.CYAN}🚗 {car} — меню автомобиля")
            print("1. Добавить ремонт")
            print("2. Добавить замену масла")
            print("3. Добавить напоминание")
            print("4. Показать историю")
            print("5. Показать статистику")
            print("6. Назад")
            choice = input("Выберите действие: ").strip()
            if choice == "1":
                self._add_repair(car)
            elif choice == "2":
                self._add_oil_change(car)
            elif choice == "3":
                self._add_reminder(car)
            elif choice == "4":
                self._show_history(car)
            elif choice == "5":
                self._show_stats(car)
            elif choice == "6":
                break
            else:
                print(Fore.RED + "Неверный выбор.")

    def _add_repair(self, car):
        date = input("Дата (ГГГГ-ММ-ДД): ")
        mileage = int(input("Пробег (км): "))
        desc = input("Описание: ")
        cost = float(input("Стоимость (руб): "))
        car.add_repair(date, mileage, desc, cost)
        self.save()
        print(Fore.GREEN + "✅ Ремонт добавлен.")

    def _add_oil_change(self, car):
        date = input("Дата (ГГГГ-ММ-ДД): ")
        mileage = int(input("Пробег (км): "))
        oil = input("Тип масла: ")
        filter_ref = input("Фильтр (артикул): ")
        cost = float(input("Стоимость (руб): "))
        interval = input("Интервал (км, по умолч. 10000): ")
        interval = int(interval) if interval.strip() else 10000
        car.add_oil_change(date, mileage, oil, filter_ref, cost, interval)
        self.save()
        print(Fore.GREEN + "✅ Замена масла добавлена.")

    def _add_reminder(self, car):
        title = input("Название: ")
        date = input("Дата (ГГГГ-ММ-ДД): ")
        mileage = int(input("Пробег (км): "))
        desc = input("Описание: ")
        priority = input("Приоритет (низкий/средний/высокий): ").lower()
        if priority not in ["низкий", "средний", "высокий"]:
            priority = "средний"
        car.add_reminder(title, date, mileage, desc, priority)
        self.save()
        print(Fore.GREEN + "✅ Напоминание добавлено.")

    def _show_history(self, car):
        print(f"\n{Fore.CYAN}📋 История для {car}")
        if car.repairs:
            print(Fore.YELLOW + "Ремонты:")
            for r in car.repairs:
                print(f"  {r['date']} | {r['mileage']} км | {r['description']} | {r['cost']} руб.")
        if car.oil_changes:
            print(Fore.YELLOW + "Замены масла:")
            for o in car.oil_changes:
                print(f"  {o['date']} | {o['mileage']} км | {o['oil_type']} | {o['filter']} | {o['cost']} руб. | след. {o['next_mileage']} км")
        if car.reminders:
            print(Fore.YELLOW + "Напоминания:")
            for r in car.reminders:
                status = "✅" if r["done"] else "⏳"
                print(f"  {status} {r['title']} | {r['date']} | {r['priority']} | {r['description']}")

    def _show_stats(self, car):
        total_repair = sum(r["cost"] for r in car.repairs)
        total_oil = sum(o["cost"] for o in car.oil_changes)
        total = total_repair + total_oil
        print(f"\n{Fore.CYAN}📊 Статистика для {car}")
        print(f"  Всего ремонтов: {len(car.repairs)}")
        print(f"  Всего замен масла: {len(car.oil_changes)}")
        print(f"  Всего напоминаний: {len(car.reminders)}")
        print(f"  Общая стоимость: {total:.2f} руб.")
        print(f"  Средняя стоимость ремонта: {total_repair / len(car.repairs) if car.repairs else 0:.2f} руб.")

def main():
    manager = CarManager()
    while True:
        print(Fore.CYAN + "\n🚗 Многомарочный дневник автомобиля (Python)")
        print("1. Выбрать/создать автомобиль")
        print("2. Удалить автомобиль")
        print("3. Показать все автомобили")
        print("4. Выход")
        choice = input("Выберите действие: ").strip()
        if choice == "1":
            manager.list_cars()
            if manager.cars:
                car_id = input("Введите номер автомобиля (или 0 для создания нового): ")
                if car_id.strip() == "0":
                    brand = input("Марка: ")
                    model = input("Модель: ")
                    year = input("Год: ")
                    vin = input("VIN: ")
                    mileage = int(input("Текущий пробег (км): "))
                    id = manager.add_car(brand, model, int(year), vin, mileage)
                    print(Fore.GREEN + f"✅ Автомобиль добавлен (ID: {id})")
                    car = manager.get_car(id)
                    manager.car_menu(car)
                elif car_id.strip().isdigit():
                    car = manager.get_car(int(car_id))
                    if car:
                        manager.car_menu(car)
                    else:
                        print(Fore.RED + "❌ Автомобиль не найден.")
            else:
                brand = input("Марка: ")
                model = input("Модель: ")
                year = input("Год: ")
                vin = input("VIN: ")
                mileage = int(input("Текущий пробег (км): "))
                id = manager.add_car(brand, model, int(year), vin, mileage)
                print(Fore.GREEN + f"✅ Автомобиль добавлен (ID: {id})")
                car = manager.get_car(id)
                manager.car_menu(car)
        elif choice == "2":
            manager.list_cars()
            car_id = int(input("Введите ID автомобиля для удаления: "))
            if manager.delete_car(car_id):
                print(Fore.GREEN + "✅ Автомобиль удалён.")
            else:
                print(Fore.RED + "❌ Автомобиль не найден.")
        elif choice == "3":
            manager.list_cars()
        elif choice == "4":
            print("До свидания!")
            break
        else:
            print(Fore.RED + "Неверный выбор.")

if __name__ == "__main__":
    main()
