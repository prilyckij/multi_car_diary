// multi_car_diary.go — Go версия

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Repair struct {
	Date        string  `json:"date"`
	Mileage     int     `json:"mileage"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

type OilChange struct {
	Date        string  `json:"date"`
	Mileage     int     `json:"mileage"`
	OilType     string  `json:"oil_type"`
	Filter      string  `json:"filter"`
	Cost        float64 `json:"cost"`
	Interval    int     `json:"interval"`
	NextMileage int     `json:"next_mileage"`
}

type Reminder struct {
	Title       string `json:"title"`
	Date        string `json:"date"`
	Mileage     int    `json:"mileage"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Done        bool   `json:"done"`
}

type Car struct {
	ID          int          `json:"id"`
	Brand       string       `json:"brand"`
	Model       string       `json:"model"`
	Year        int          `json:"year"`
	Vin         string       `json:"vin"`
	Mileage     int          `json:"mileage"`
	Repairs     []Repair     `json:"repairs"`
	OilChanges  []OilChange  `json:"oil_changes"`
	Reminders   []Reminder   `json:"reminders"`
}

type Manager struct {
	Cars []Car `json:"cars"`
	file string
}

func NewManager(file string) *Manager {
	m := &Manager{file: file}
	m.load()
	return m
}

func (m *Manager) load() {
	data, err := os.ReadFile(m.file)
	if err != nil {
		m.Cars = []Car{}
		return
	}
	json.Unmarshal(data, &m.Cars)
}

func (m *Manager) save() {
	data, _ := json.MarshalIndent(m.Cars, "", "  ")
	os.WriteFile(m.file, data, 0644)
}

func (m *Manager) addCar(brand, model string, year int, vin string, mileage int) int {
	id := len(m.Cars) + 1
	m.Cars = append(m.Cars, Car{
		ID:          id,
		Brand:       brand,
		Model:       model,
		Year:        year,
		Vin:         vin,
		Mileage:     mileage,
		Repairs:     []Repair{},
		OilChanges:  []OilChange{},
		Reminders:   []Reminder{},
	})
	m.save()
	return id
}

func (m *Manager) deleteCar(id int) bool {
	for i, c := range m.Cars {
		if c.ID == id {
			m.Cars = append(m.Cars[:i], m.Cars[i+1:]...)
			m.save()
			return true
		}
	}
	return false
}

func (m *Manager) getCar(id int) *Car {
	for i := range m.Cars {
		if m.Cars[i].ID == id {
			return &m.Cars[i]
		}
	}
	return nil
}

func (m *Manager) listCars() {
	if len(m.Cars) == 0 {
		fmt.Println("\u001B[33mНет автомобилей.\u001B[0m")
		return
	}
	fmt.Printf("\u001B[36m%-4s %-15s %-15s %-6s %-10s\u001B[0m\n", "ID", "Марка", "Модель", "Год", "Пробег")
	fmt.Println(strings.Repeat("-", 55))
	for _, c := range m.Cars {
		fmt.Printf("%-4d %-15s %-15s %-6d %-10d\n", c.ID, c.Brand, c.Model, c.Year, c.Mileage)
	}
}

func (m *Manager) carMenu(car *Car) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n\u001B[36m🚗 %s %s %d (Пробег: %d км) — меню автомобиля\u001B[0m\n", car.Brand, car.Model, car.Year, car.Mileage)
		fmt.Println("1. Добавить ремонт")
		fmt.Println("2. Добавить замену масла")
		fmt.Println("3. Добавить напоминание")
		fmt.Println("4. Показать историю")
		fmt.Println("5. Показать статистику")
		fmt.Println("6. Назад")
		fmt.Print("Выберите действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			m.addRepair(car, reader)
		case "2":
			m.addOilChange(car, reader)
		case "3":
			m.addReminder(car, reader)
		case "4":
			m.showHistory(car)
		case "5":
			m.showStats(car)
		case "6":
			return
		default:
			fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
		}
	}
}

func (m *Manager) addRepair(car *Car, reader *bufio.Reader) {
	fmt.Print("Дата (ГГГГ-ММ-ДД): ")
	date, _ := reader.ReadString('\n')
	date = strings.TrimSpace(date)
	fmt.Print("Пробег (км): ")
	mileageStr, _ := reader.ReadString('\n')
	mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
	fmt.Print("Описание: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	fmt.Print("Стоимость (руб): ")
	costStr, _ := reader.ReadString('\n')
	cost, _ := strconv.ParseFloat(strings.TrimSpace(costStr), 64)

	car.Repairs = append(car.Repairs, Repair{Date: date, Mileage: mileage, Description: desc, Cost: cost})
	if mileage > car.Mileage {
		car.Mileage = mileage
	}
	m.save()
	fmt.Println("\u001B[32m✅ Ремонт добавлен.\u001B[0m")
}

func (m *Manager) addOilChange(car *Car, reader *bufio.Reader) {
	fmt.Print("Дата (ГГГГ-ММ-ДД): ")
	date, _ := reader.ReadString('\n')
	date = strings.TrimSpace(date)
	fmt.Print("Пробег (км): ")
	mileageStr, _ := reader.ReadString('\n')
	mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
	fmt.Print("Тип масла: ")
	oil, _ := reader.ReadString('\n')
	oil = strings.TrimSpace(oil)
	fmt.Print("Фильтр (артикул): ")
	filter, _ := reader.ReadString('\n')
	filter = strings.TrimSpace(filter)
	fmt.Print("Стоимость (руб): ")
	costStr, _ := reader.ReadString('\n')
	cost, _ := strconv.ParseFloat(strings.TrimSpace(costStr), 64)
	fmt.Print("Интервал (км, по умолч. 10000): ")
	intervalStr, _ := reader.ReadString('\n')
	interval := 10000
	if intervalStr = strings.TrimSpace(intervalStr); intervalStr != "" {
		interval, _ = strconv.Atoi(intervalStr)
	}

	nextMileage := mileage + interval
	car.OilChanges = append(car.OilChanges, OilChange{
		Date: date, Mileage: mileage, OilType: oil, Filter: filter,
		Cost: cost, Interval: interval, NextMileage: nextMileage,
	})
	if mileage > car.Mileage {
		car.Mileage = mileage
	}
	m.save()
	fmt.Println("\u001B[32m✅ Замена масла добавлена.\u001B[0m")
}

func (m *Manager) addReminder(car *Car, reader *bufio.Reader) {
	fmt.Print("Название: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	fmt.Print("Дата (ГГГГ-ММ-ДД): ")
	date, _ := reader.ReadString('\n')
	date = strings.TrimSpace(date)
	fmt.Print("Пробег (км): ")
	mileageStr, _ := reader.ReadString('\n')
	mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
	fmt.Print("Описание: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	fmt.Print("Приоритет (низкий/средний/высокий): ")
	priority, _ := reader.ReadString('\n')
	priority = strings.TrimSpace(strings.ToLower(priority))
	if priority != "низкий" && priority != "средний" && priority != "высокий" {
		priority = "средний"
	}
	car.Reminders = append(car.Reminders, Reminder{
		Title: title, Date: date, Mileage: mileage, Description: desc, Priority: priority, Done: false,
	})
	m.save()
	fmt.Println("\u001B[32m✅ Напоминание добавлено.\u001B[0m")
}

func (m *Manager) showHistory(car *Car) {
	fmt.Printf("\n\u001B[36m📋 История для %s %s %d\u001B[0m\n", car.Brand, car.Model, car.Year)
	if len(car.Repairs) > 0 {
		fmt.Println("\u001B[33mРемонты:\u001B[0m")
		for _, r := range car.Repairs {
			fmt.Printf("  %s | %d км | %s | %.2f руб.\n", r.Date, r.Mileage, r.Description, r.Cost)
		}
	}
	if len(car.OilChanges) > 0 {
		fmt.Println("\u001B[33mЗамены масла:\u001B[0m")
		for _, o := range car.OilChanges {
			fmt.Printf("  %s | %d км | %s | %s | %.2f руб. | след. %d км\n", o.Date, o.Mileage, o.OilType, o.Filter, o.Cost, o.NextMileage)
		}
	}
	if len(car.Reminders) > 0 {
		fmt.Println("\u001B[33mНапоминания:\u001B[0m")
		for _, r := range car.Reminders {
			status := "⏳"
			if r.Done {
				status = "✅"
			}
			fmt.Printf("  %s %s | %s | %s | %s\n", status, r.Title, r.Date, r.Priority, r.Description)
		}
	}
}

func (m *Manager) showStats(car *Car) {
	totalRepair := 0.0
	for _, r := range car.Repairs {
		totalRepair += r.Cost
	}
	totalOil := 0.0
	for _, o := range car.OilChanges {
		totalOil += o.Cost
	}
	total := totalRepair + totalOil
	fmt.Printf("\n\u001B[36m📊 Статистика для %s %s %d\u001B[0m\n", car.Brand, car.Model, car.Year)
	fmt.Printf("  Всего ремонтов: %d\n", len(car.Repairs))
	fmt.Printf("  Всего замен масла: %d\n", len(car.OilChanges))
	fmt.Printf("  Всего напоминаний: %d\n", len(car.Reminders))
	fmt.Printf("  Общая стоимость: %.2f руб.\n", total)
	if len(car.Repairs) > 0 {
		fmt.Printf("  Средняя стоимость ремонта: %.2f руб.\n", totalRepair/float64(len(car.Repairs)))
	}
}

func main() {
	manager := NewManager("cars.json")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n\u001B[36m🚗 Многомарочный дневник автомобиля (Go)\u001B[0m")
		fmt.Println("1. Выбрать/создать автомобиль")
		fmt.Println("2. Удалить автомобиль")
		fmt.Println("3. Показать все автомобили")
		fmt.Println("4. Выход")
		fmt.Print("Выберите действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			manager.listCars()
			if len(manager.Cars) > 0 {
				fmt.Print("Введите номер автомобиля (или 0 для создания нового): ")
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr)
				if idStr == "0" {
					fmt.Print("Марка: ")
					brand, _ := reader.ReadString('\n')
					brand = strings.TrimSpace(brand)
					fmt.Print("Модель: ")
					model, _ := reader.ReadString('\n')
					model = strings.TrimSpace(model)
					fmt.Print("Год: ")
					yearStr, _ := reader.ReadString('\n')
					year, _ := strconv.Atoi(strings.TrimSpace(yearStr))
					fmt.Print("VIN: ")
					vin, _ := reader.ReadString('\n')
					vin = strings.TrimSpace(vin)
					fmt.Print("Текущий пробег (км): ")
					mileageStr, _ := reader.ReadString('\n')
					mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
					id := manager.addCar(brand, model, year, vin, mileage)
					fmt.Printf("\u001B[32m✅ Автомобиль добавлен (ID: %d)\u001B[0m\n", id)
					car := manager.getCar(id)
					if car != nil {
						manager.carMenu(car)
					}
				} else if idStr != "" {
					id, _ := strconv.Atoi(idStr)
					car := manager.getCar(id)
					if car != nil {
						manager.carMenu(car)
					} else {
						fmt.Println("\u001B[31m❌ Автомобиль не найден.\u001B[0m")
					}
				}
			} else {
				fmt.Print("Марка: ")
				brand, _ := reader.ReadString('\n')
				brand = strings.TrimSpace(brand)
				fmt.Print("Модель: ")
				model, _ := reader.ReadString('\n')
				model = strings.TrimSpace(model)
				fmt.Print("Год: ")
				yearStr, _ := reader.ReadString('\n')
				year, _ := strconv.Atoi(strings.TrimSpace(yearStr))
				fmt.Print("VIN: ")
				vin, _ := reader.ReadString('\n')
				vin = strings.TrimSpace(vin)
				fmt.Print("Текущий пробег (км): ")
				mileageStr, _ := reader.ReadString('\n')
				mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
				id := manager.addCar(brand, model, year, vin, mileage)
				fmt.Printf("\u001B[32m✅ Автомобиль добавлен (ID: %d)\u001B[0m\n", id)
				car := manager.getCar(id)
				if car != nil {
					manager.carMenu(car)
				}
			}
		case "2":
			manager.listCars()
			fmt.Print("Введите ID автомобиля для удаления: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			if manager.deleteCar(id) {
				fmt.Println("\u001B[32m✅ Автомобиль удалён.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Автомобиль не найден.\u001B[0m")
			}
		case "3":
			manager.listCars()
		case "4":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
		}
	}
}
