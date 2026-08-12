// multi_car_diary.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;

class Repair
{
    public string Date { get; set; }
    public int Mileage { get; set; }
    public string Description { get; set; }
    public double Cost { get; set; }
}

class OilChange
{
    public string Date { get; set; }
    public int Mileage { get; set; }
    public string OilType { get; set; }
    public string Filter { get; set; }
    public double Cost { get; set; }
    public int Interval { get; set; }
    public int NextMileage { get; set; }
}

class Reminder
{
    public string Title { get; set; }
    public string Date { get; set; }
    public int Mileage { get; set; }
    public string Description { get; set; }
    public string Priority { get; set; }
    public bool Done { get; set; }
}

class Car
{
    public int Id { get; set; }
    public string Brand { get; set; }
    public string Model { get; set; }
    public int Year { get; set; }
    public string Vin { get; set; }
    public int Mileage { get; set; }
    public List<Repair> Repairs { get; set; } = new List<Repair>();
    public List<OilChange> OilChanges { get; set; } = new List<OilChange>();
    public List<Reminder> Reminders { get; set; } = new List<Reminder>();
}

class Program
{
    private static List<Car> cars = new List<Car>();
    private const string DataFile = "cars.json";

    static void Main()
    {
        Load();
        while (true)
        {
            Console.WriteLine("\n\u001B[36m🚗 Многомарочный дневник автомобиля (C#)\u001B[0m");
            Console.WriteLine("1. Выбрать/создать автомобиль");
            Console.WriteLine("2. Удалить автомобиль");
            Console.WriteLine("3. Показать все автомобили");
            Console.WriteLine("4. Выход");
            Console.Write("Выберите действие: ");
            string choice = Console.ReadLine();
            switch (choice)
            {
                case "1": SelectOrCreateCar(); break;
                case "2": DeleteCar(); break;
                case "3": ListCars(); break;
                case "4": Console.WriteLine("До свидания!"); return;
                default: Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m"); break;
            }
        }
    }

    static void Load()
    {
        if (File.Exists(DataFile))
        {
            try
            {
                string json = File.ReadAllText(DataFile);
                cars = JsonSerializer.Deserialize<List<Car>>(json) ?? new List<Car>();
            }
            catch { cars = new List<Car>(); }
        }
    }

    static void Save()
    {
        string json = JsonSerializer.Serialize(cars, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(DataFile, json);
    }

    static Car GetCar(int id) => cars.FirstOrDefault(c => c.Id == id);

    static void ListCars()
    {
        if (cars.Count == 0)
        {
            Console.WriteLine("\u001B[33mНет автомобилей.\u001B[0m");
            return;
        }
        Console.WriteLine($"\u001B[36m{"ID",-4} {"Марка",-15} {"Модель",-15} {"Год",-6} {"Пробег",-10}\u001B[0m");
        Console.WriteLine(new string('-', 55));
        foreach (var c in cars)
            Console.WriteLine($"{c.Id,-4} {c.Brand,-15} {c.Model,-15} {c.Year,-6} {c.Mileage,-10}");
    }

    static void SelectOrCreateCar()
    {
        ListCars();
        if (cars.Count > 0)
        {
            Console.Write("Введите номер автомобиля (или 0 для создания нового): ");
            string input = Console.ReadLine();
            if (input == "0")
            {
                CreateNewCar();
            }
            else if (int.TryParse(input, out int id))
            {
                var car = GetCar(id);
                if (car != null)
                    CarMenu(car);
                else
                    Console.WriteLine("\u001B[31m❌ Автомобиль не найден.\u001B[0m");
            }
        }
        else
        {
            CreateNewCar();
        }
    }

    static void CreateNewCar()
    {
        Console.Write("Марка: ");
        string brand = Console.ReadLine();
        Console.Write("Модель: ");
        string model = Console.ReadLine();
        Console.Write("Год: ");
        int year = int.Parse(Console.ReadLine());
        Console.Write("VIN: ");
        string vin = Console.ReadLine();
        Console.Write("Текущий пробег (км): ");
        int mileage = int.Parse(Console.ReadLine());
        int id = cars.Count + 1;
        cars.Add(new Car { Id = id, Brand = brand, Model = model, Year = year, Vin = vin, Mileage = mileage });
        Save();
        Console.WriteLine($"\u001B[32m✅ Автомобиль добавлен (ID: {id})\u001B[0m");
        var car = GetCar(id);
        if (car != null) CarMenu(car);
    }

    static void DeleteCar()
    {
        ListCars();
        Console.Write("Введите ID автомобиля для удаления: ");
        int id = int.Parse(Console.ReadLine());
        var car = GetCar(id);
        if (car != null)
        {
            cars.Remove(car);
            Save();
            Console.WriteLine("\u001B[32m✅ Автомобиль удалён.\u001B[0m");
        }
        else
            Console.WriteLine("\u001B[31m❌ Автомобиль не найден.\u001B[0m");
    }

    static void CarMenu(Car car)
    {
        while (true)
        {
            Console.WriteLine($"\n\u001B[36m🚗 {car.Brand} {car.Model} {car.Year} (Пробег: {car.Mileage} км) — меню автомобиля\u001B[0m");
            Console.WriteLine("1. Добавить ремонт");
            Console.WriteLine("2. Добавить замену масла");
            Console.WriteLine("3. Добавить напоминание");
            Console.WriteLine("4. Показать историю");
            Console.WriteLine("5. Показать статистику");
            Console.WriteLine("6. Назад");
            Console.Write("Выберите действие: ");
            string choice = Console.ReadLine();
            switch (choice)
            {
                case "1": AddRepair(car); break;
                case "2": AddOilChange(car); break;
                case "3": AddReminder(car); break;
                case "4": ShowHistory(car); break;
                case "5": ShowStats(car); break;
                case "6": return;
                default: Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m"); break;
            }
        }
    }

    static void AddRepair(Car car)
    {
        Console.Write("Дата (ГГГГ-ММ-ДД): ");
        string date = Console.ReadLine();
        Console.Write("Пробег (км): ");
        int mileage = int.Parse(Console.ReadLine());
        Console.Write("Описание: ");
        string desc = Console.ReadLine();
        Console.Write("Стоимость (руб): ");
        double cost = double.Parse(Console.ReadLine());
        car.Repairs.Add(new Repair { Date = date, Mileage = mileage, Description = desc, Cost = cost });
        if (mileage > car.Mileage) car.Mileage = mileage;
        Save();
        Console.WriteLine("\u001B[32m✅ Ремонт добавлен.\u001B[0m");
    }

    static void AddOilChange(Car car)
    {
        Console.Write("Дата (ГГГГ-ММ-ДД): ");
        string date = Console.ReadLine();
        Console.Write("Пробег (км): ");
        int mileage = int.Parse(Console.ReadLine());
        Console.Write("Тип масла: ");
        string oil = Console.ReadLine();
        Console.Write("Фильтр (артикул): ");
        string filter = Console.ReadLine();
        Console.Write("Стоимость (руб): ");
        double cost = double.Parse(Console.ReadLine());
        Console.Write("Интервал (км, по умолч. 10000): ");
        string intervalStr = Console.ReadLine();
        int interval = string.IsNullOrEmpty(intervalStr) ? 10000 : int.Parse(intervalStr);
        car.OilChanges.Add(new OilChange { Date = date, Mileage = mileage, OilType = oil, Filter = filter, Cost = cost, Interval = interval, NextMileage = mileage + interval });
        if (mileage > car.Mileage) car.Mileage = mileage;
        Save();
        Console.WriteLine("\u001B[32m✅ Замена масла добавлена.\u001B[0m");
    }

    static void AddReminder(Car car)
    {
        Console.Write("Название: ");
        string title = Console.ReadLine();
        Console.Write("Дата (ГГГГ-ММ-ДД): ");
        string date = Console.ReadLine();
        Console.Write("Пробег (км): ");
        int mileage = int.Parse(Console.ReadLine());
        Console.Write("Описание: ");
        string desc = Console.ReadLine();
        Console.Write("Приоритет (низкий/средний/высокий): ");
        string priority = Console.ReadLine().ToLower();
        if (priority != "низкий" && priority != "средний" && priority != "высокий") priority = "средний";
        car.Reminders.Add(new Reminder { Title = title, Date = date, Mileage = mileage, Description = desc, Priority = priority, Done = false });
        Save();
        Console.WriteLine("\u001B[32m✅ Напоминание добавлено.\u001B[0m");
    }

    static void ShowHistory(Car car)
    {
        Console.WriteLine($"\n\u001B[36m📋 История для {car.Brand} {car.Model} {car.Year}\u001B[0m");
        if (car.Repairs.Count > 0)
        {
            Console.WriteLine("\u001B[33mРемонты:\u001B[0m");
            foreach (var r in car.Repairs)
                Console.WriteLine($"  {r.Date} | {r.Mileage} км | {r.Description} | {r.Cost} руб.");
        }
        if (car.OilChanges.Count > 0)
        {
            Console.WriteLine("\u001B[33mЗамены масла:\u001B[0m");
            foreach (var o in car.OilChanges)
                Console.WriteLine($"  {o.Date} | {o.Mileage} км | {o.OilType} | {o.Filter} | {o.Cost} руб. | след. {o.NextMileage} км");
        }
        if (car.Reminders.Count > 0)
        {
            Console.WriteLine("\u001B[33mНапоминания:\u001B[0m");
            foreach (var r in car.Reminders)
                Console.WriteLine($"  {(r.Done ? "✅" : "⏳")} {r.Title} | {r.Date} | {r.Priority} | {r.Description}");
        }
    }

    static void ShowStats(Car car)
    {
        double totalRepair = car.Repairs.Sum(r => r.Cost);
        double totalOil = car.OilChanges.Sum(o => o.Cost);
        double total = totalRepair + totalOil;
        Console.WriteLine($"\n\u001B[36m📊 Статистика для {car.Brand} {car.Model} {car.Year}\u001B[0m");
        Console.WriteLine($"  Всего ремонтов: {car.Repairs.Count}");
        Console.WriteLine($"  Всего замен масла: {car.OilChanges.Count}");
        Console.WriteLine($"  Всего напоминаний: {car.Reminders.Count}");
        Console.WriteLine($"  Общая стоимость: {total:F2} руб.");
        if (car.Repairs.Count > 0)
            Console.WriteLine($"  Средняя стоимость ремонта: {totalRepair / car.Repairs.Count:F2} руб.");
    }
}
