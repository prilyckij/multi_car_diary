// multi_car_diary.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.time.*;
import java.time.format.*;

class Repair {
    String date;
    int mileage;
    String description;
    double cost;

    Repair(String date, int mileage, String description, double cost) {
        this.date = date;
        this.mileage = mileage;
        this.description = description;
        this.cost = cost;
    }

    String toJson() {
        return String.format("{\"date\":\"%s\",\"mileage\":%d,\"description\":\"%s\",\"cost\":%.2f}",
                date, mileage, description, cost);
    }
}

class OilChange {
    String date;
    int mileage;
    String oilType;
    String filter;
    double cost;
    int interval;
    int nextMileage;

    OilChange(String date, int mileage, String oilType, String filter, double cost, int interval, int nextMileage) {
        this.date = date;
        this.mileage = mileage;
        this.oilType = oilType;
        this.filter = filter;
        this.cost = cost;
        this.interval = interval;
        this.nextMileage = nextMileage;
    }

    String toJson() {
        return String.format("{\"date\":\"%s\",\"mileage\":%d,\"oil_type\":\"%s\",\"filter\":\"%s\",\"cost\":%.2f,\"interval\":%d,\"next_mileage\":%d}",
                date, mileage, oilType, filter, cost, interval, nextMileage);
    }
}

class Reminder {
    String title;
    String date;
    int mileage;
    String description;
    String priority;
    boolean done;

    Reminder(String title, String date, int mileage, String description, String priority, boolean done) {
        this.title = title;
        this.date = date;
        this.mileage = mileage;
        this.description = description;
        this.priority = priority;
        this.done = done;
    }

    String toJson() {
        return String.format("{\"title\":\"%s\",\"date\":\"%s\",\"mileage\":%d,\"description\":\"%s\",\"priority\":\"%s\",\"done\":%b}",
                title, date, mileage, description, priority, done);
    }
}

class Car {
    int id;
    String brand;
    String model;
    int year;
    String vin;
    int mileage;
    List<Repair> repairs;
    List<OilChange> oilChanges;
    List<Reminder> reminders;

    Car(int id, String brand, String model, int year, String vin, int mileage) {
        this.id = id;
        this.brand = brand;
        this.model = model;
        this.year = year;
        this.vin = vin;
        this.mileage = mileage;
        this.repairs = new ArrayList<>();
        this.oilChanges = new ArrayList<>();
        this.reminders = new ArrayList<>();
    }

    String toJson() {
        StringBuilder sb = new StringBuilder();
        sb.append("{\"id\":").append(id).append(",\"brand\":\"").append(brand).append("\",\"model\":\"").append(model)
          .append("\",\"year\":").append(year).append(",\"vin\":\"").append(vin).append("\",\"mileage\":").append(mileage)
          .append(",\"repairs\":[");
        for (int i = 0; i < repairs.size(); i++) {
            sb.append(repairs.get(i).toJson());
            if (i < repairs.size() - 1) sb.append(",");
        }
        sb.append("],\"oil_changes\":[");
        for (int i = 0; i < oilChanges.size(); i++) {
            sb.append(oilChanges.get(i).toJson());
            if (i < oilChanges.size() - 1) sb.append(",");
        }
        sb.append("],\"reminders\":[");
        for (int i = 0; i < reminders.size(); i++) {
            sb.append(reminders.get(i).toJson());
            if (i < reminders.size() - 1) sb.append(",");
        }
        sb.append("]}");
        return sb.toString();
    }
}

public class multi_car_diary {
    private static final String DATA_FILE = "cars.json";
    private static List<Car> cars = new ArrayList<>();
    private static Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) {
        load();
        while (true) {
            System.out.println("\n\u001B[36m🚗 Многомарочный дневник автомобиля (Java)\u001B[0m");
            System.out.println("1. Выбрать/создать автомобиль");
            System.out.println("2. Удалить автомобиль");
            System.out.println("3. Показать все автомобили");
            System.out.println("4. Выход");
            System.out.print("Выберите действие: ");
            String choice = scanner.nextLine();
            switch (choice) {
                case "1": selectOrCreateCar(); break;
                case "2": deleteCar(); break;
                case "3": listCars(); break;
                case "4": System.out.println("До свидания!"); return;
                default: System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
            }
        }
    }

    private static void load() {
        try {
            String content = new String(Files.readAllBytes(Paths.get(DATA_FILE)));
            // Упрощённый парсинг, в реальности нужен JSON-парсер
            cars = new ArrayList<>();
        } catch (IOException e) {
            cars = new ArrayList<>();
        }
    }

    private static void save() {
        try {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < cars.size(); i++) {
                sb.append(cars.get(i).toJson());
                if (i < cars.size() - 1) sb.append(",");
            }
            sb.append("]");
            Files.write(Paths.get(DATA_FILE), sb.toString().getBytes());
        } catch (IOException e) {
            System.out.println("Ошибка сохранения.");
        }
    }

    private static void listCars() {
        if (cars.isEmpty()) {
            System.out.println("\u001B[33mНет автомобилей.\u001B[0m");
            return;
        }
        System.out.printf("\u001B[36m%-4s %-15s %-15s %-6s %-10s\u001B[0m\n", "ID", "Марка", "Модель", "Год", "Пробег");
        System.out.println("-".repeat(55));
        for (Car c : cars) {
            System.out.printf("%-4d %-15s %-15s %-6d %-10d\n", c.id, c.brand, c.model, c.year, c.mileage);
        }
    }

    private static Car getCar(int id) {
        for (Car c : cars) {
            if (c.id == id) return c;
        }
        return null;
    }

    private static int addCar(String brand, String model, int year, String vin, int mileage) {
        int id = cars.size() + 1;
        cars.add(new Car(id, brand, model, year, vin, mileage));
        save();
        return id;
    }

    private static void deleteCar() {
        listCars();
        System.out.print("Введите ID автомобиля для удаления: ");
        int id = Integer.parseInt(scanner.nextLine());
        boolean removed = cars.removeIf(c -> c.id == id);
        if (removed) {
            save();
            System.out.println("\u001B[32m✅ Автомобиль удалён.\u001B[0m");
        } else {
            System.out.println("\u001B[31m❌ Автомобиль не найден.\u001B[0m");
        }
    }

    private static void selectOrCreateCar() {
        listCars();
        if (!cars.isEmpty()) {
            System.out.print("Введите номер автомобиля (или 0 для создания нового): ");
            String input = scanner.nextLine();
            if (input.equals("0")) {
                createNewCar();
            } else {
                int id = Integer.parseInt(input);
                Car car = getCar(id);
                if (car != null) {
                    carMenu(car);
                } else {
                    System.out.println("\u001B[31m❌ Автомобиль не найден.\u001B[0m");
                }
            }
        } else {
            createNewCar();
        }
    }

    private static void createNewCar() {
        System.out.print("Марка: ");
        String brand = scanner.nextLine();
        System.out.print("Модель: ");
        String model = scanner.nextLine();
        System.out.print("Год: ");
        int year = Integer.parseInt(scanner.nextLine());
        System.out.print("VIN: ");
        String vin = scanner.nextLine();
        System.out.print("Текущий пробег (км): ");
        int mileage = Integer.parseInt(scanner.nextLine());
        int id = addCar(brand, model, year, vin, mileage);
        System.out.println("\u001B[32m✅ Автомобиль добавлен (ID: " + id + ")\u001B[0m");
        Car car = getCar(id);
        if (car != null) carMenu(car);
    }

    private static void carMenu(Car car) {
        while (true) {
            System.out.println("\n\u001B[36m🚗 " + car.brand + " " + car.model + " " + car.year + " (Пробег: " + car.mileage + " км) — меню автомобиля\u001B[0m");
            System.out.println("1. Добавить ремонт");
            System.out.println("2. Добавить замену масла");
            System.out.println("3. Добавить напоминание");
            System.out.println("4. Показать историю");
            System.out.println("5. Показать статистику");
            System.out.println("6. Назад");
            System.out.print("Выберите действие: ");
            String choice = scanner.nextLine();
            switch (choice) {
                case "1": addRepair(car); break;
                case "2": addOilChange(car); break;
                case "3": addReminder(car); break;
                case "4": showHistory(car); break;
                case "5": showStats(car); break;
                case "6": return;
                default: System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
            }
        }
    }

    private static void addRepair(Car car) {
        System.out.print("Дата (ГГГГ-ММ-ДД): ");
        String date = scanner.nextLine();
        System.out.print("Пробег (км): ");
        int mileage = Integer.parseInt(scanner.nextLine());
        System.out.print("Описание: ");
        String desc = scanner.nextLine();
        System.out.print("Стоимость (руб): ");
        double cost = Double.parseDouble(scanner.nextLine());
        car.repairs.add(new Repair(date, mileage, desc, cost));
        if (mileage > car.mileage) car.mileage = mileage;
        save();
        System.out.println("\u001B[32m✅ Ремонт добавлен.\u001B[0m");
    }

    private static void addOilChange(Car car) {
        System.out.print("Дата (ГГГГ-ММ-ДД): ");
        String date = scanner.nextLine();
        System.out.print("Пробег (км): ");
        int mileage = Integer.parseInt(scanner.nextLine());
        System.out.print("Тип масла: ");
        String oil = scanner.nextLine();
        System.out.print("Фильтр (артикул): ");
        String filter = scanner.nextLine();
        System.out.print("Стоимость (руб): ");
        double cost = Double.parseDouble(scanner.nextLine());
        System.out.print("Интервал (км, по умолч. 10000): ");
        String intervalStr = scanner.nextLine();
        int interval = intervalStr.isEmpty() ? 10000 : Integer.parseInt(intervalStr);
        car.oilChanges.add(new OilChange(date, mileage, oil, filter, cost, interval, mileage + interval));
        if (mileage > car.mileage) car.mileage = mileage;
        save();
        System.out.println("\u001B[32m✅ Замена масла добавлена.\u001B[0m");
    }

    private static void addReminder(Car car) {
        System.out.print("Название: ");
        String title = scanner.nextLine();
        System.out.print("Дата (ГГГГ-ММ-ДД): ");
        String date = scanner.nextLine();
        System.out.print("Пробег (км): ");
        int mileage = Integer.parseInt(scanner.nextLine());
        System.out.print("Описание: ");
        String desc = scanner.nextLine();
        System.out.print("Приоритет (низкий/средний/высокий): ");
        String priority = scanner.nextLine().toLowerCase();
        if (!priority.equals("низкий") && !priority.equals("средний") && !priority.equals("высокий")) {
            priority = "средний";
        }
        car.reminders.add(new Reminder(title, date, mileage, desc, priority, false));
        save();
        System.out.println("\u001B[32m✅ Напоминание добавлено.\u001B[0m");
    }

    private static void showHistory(Car car) {
        System.out.println("\n\u001B[36m📋 История для " + car.brand + " " + car.model + " " + car.year + "\u001B[0m");
        if (!car.repairs.isEmpty()) {
            System.out.println("\u001B[33mРемонты:\u001B[0m");
            for (Repair r : car.repairs) {
                System.out.printf("  %s | %d км | %s | %.2f руб.\n", r.date, r.mileage, r.description, r.cost);
            }
        }
        if (!car.oilChanges.isEmpty()) {
            System.out.println("\u001B[33mЗамены масла:\u001B[0m");
            for (OilChange o : car.oilChanges) {
                System.out.printf("  %s | %d км | %s | %s | %.2f руб. | след. %d км\n", o.date, o.mileage, o.oilType, o.filter, o.cost, o.nextMileage);
            }
        }
        if (!car.reminders.isEmpty()) {
            System.out.println("\u001B[33mНапоминания:\u001B[0m");
            for (Reminder r : car.reminders) {
                String status = r.done ? "✅" : "⏳";
                System.out.printf("  %s %s | %s | %s | %s\n", status, r.title, r.date, r.priority, r.description);
            }
        }
    }

    private static void showStats(Car car) {
        double totalRepair = car.repairs.stream().mapToDouble(r -> r.cost).sum();
        double totalOil = car.oilChanges.stream().mapToDouble(o -> o.cost).sum();
        double total = totalRepair + totalOil;
        System.out.println("\n\u001B[36m📊 Статистика для " + car.brand + " " + car.model + " " + car.year + "\u001B[0m");
        System.out.println("  Всего ремонтов: " + car.repairs.size());
        System.out.println("  Всего замен масла: " + car.oilChanges.size());
        System.out.println("  Всего напоминаний: " + car.reminders.size());
        System.out.printf("  Общая стоимость: %.2f руб.\n", total);
        if (!car.repairs.isEmpty()) {
            System.out.printf("  Средняя стоимость ремонта: %.2f руб.\n", totalRepair / car.repairs.size());
        }
    }
}
