# multi_car_diary.rb — Ruby версия

require 'json'
require 'date'

DATA_FILE = 'cars.json'

class Car
  attr_accessor :id, :brand, :model, :year, :vin, :mileage, :repairs, :oil_changes, :reminders

  def initialize(id, brand, model, year, vin, mileage, repairs = [], oil_changes = [], reminders = [])
    @id = id
    @brand = brand
    @model = model
    @year = year
    @vin = vin
    @mileage = mileage
    @repairs = repairs
    @oil_changes = oil_changes
    @reminders = reminders
  end

  def to_h
    {
      id: @id, brand: @brand, model: @model, year: @year,
      vin: @vin, mileage: @mileage,
      repairs: @repairs, oil_changes: @oil_changes, reminders: @reminders
    }
  end

  def self.from_h(data)
    new(data[:id], data[:brand], data[:model], data[:year],
        data[:vin], data[:mileage], data[:repairs] || [], data[:oil_changes] || [], data[:reminders] || [])
  end

  def add_repair(date, mileage, description, cost)
    @repairs << { date: date, mileage: mileage, description: description, cost: cost }
    @mileage = [@mileage, mileage].max
  end

  def add_oil_change(date, mileage, oil_type, filter_ref, cost, interval = 10000)
    @oil_changes << {
      date: date, mileage: mileage, oil_type: oil_type, filter: filter_ref,
      cost: cost, interval: interval, next_mileage: mileage + interval
    }
    @mileage = [@mileage, mileage].max
  end

  def add_reminder(title, date, mileage, description, priority)
    @reminders << {
      title: title, date: date, mileage: mileage,
      description: description, priority: priority, done: false
    }
  end

  def total_cost
    repairs.sum { |r| r[:cost] } + oil_changes.sum { |o| o[:cost] }
  end

  def to_s
    "#{@brand} #{@model} #{@year} (Пробег: #{@mileage} км)"
  end
end

class CarManager
  attr_reader :cars

  def initialize
    @cars = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      begin
        data = JSON.parse(File.read(DATA_FILE), symbolize_names: true)
        @cars = data.map { |c| Car.from_h(c) }
      rescue
        @cars = []
      end
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate(@cars.map(&:to_h)))
  end

  def add_car(brand, model, year, vin, mileage)
    id = @cars.size + 1
    @cars << Car.new(id, brand, model, year, vin, mileage)
    save
    id
  end

  def delete_car(id)
    found = @cars.find { |c| c.id == id }
    if found
      @cars.delete(found)
      save
      true
    else
      false
    end
  end

  def get_car(id)
    @cars.find { |c| c.id == id }
  end

  def list_cars
    if @cars.empty?
      puts "\e[33mНет автомобилей.\e[0m"
      return
    end
    printf "\e[36m%-4s %-15s %-15s %-6s %-10s\e[0m\n", "ID", "Марка", "Модель", "Год", "Пробег"
    puts "-" * 55
    @cars.each do |c|
      puts "%-4d %-15s %-15s %-6d %-10d" % [c.id, c.brand, c.model, c.year, c.mileage]
    end
  end

  def car_menu(car)
    loop do
      puts "\n\e[36m🚗 #{car} — меню автомобиля\e[0m"
      puts "1. Добавить ремонт"
      puts "2. Добавить замену масла"
      puts "3. Добавить напоминание"
      puts "4. Показать историю"
      puts "5. Показать статистику"
      puts "6. Назад"
      print "Выберите действие: "
      choice = gets.chomp
      case choice
      when "1" then add_repair(car)
      when "2" then add_oil_change(car)
      when "3" then add_reminder(car)
      when "4" then show_history(car)
      when "5" then show_stats(car)
      when "6" then break
      else puts "\e[31mНеверный выбор.\e[0m"
      end
    end
  end

  def add_repair(car)
    print "Дата (ГГГГ-ММ-ДД): "
    date = gets.chomp
    print "Пробег (км): "
    mileage = gets.chomp.to_i
    print "Описание: "
    desc = gets.chomp
    print "Стоимость (руб): "
    cost = gets.chomp.to_f
    car.add_repair(date, mileage, desc, cost)
    save
    puts "\e[32m✅ Ремонт добавлен.\e[0m"
  end

  def add_oil_change(car)
    print "Дата (ГГГГ-ММ-ДД): "
    date = gets.chomp
    print "Пробег (км): "
    mileage = gets.chomp.to_i
    print "Тип масла: "
    oil = gets.chomp
    print "Фильтр (артикул): "
    filter = gets.chomp
    print "Стоимость (руб): "
    cost = gets.chomp.to_f
    print "Интервал (км, по умолч. 10000): "
    interval_str = gets.chomp
    interval = interval_str.empty? ? 10000 : interval_str.to_i
    car.add_oil_change(date, mileage, oil, filter, cost, interval)
    save
    puts "\e[32m✅ Замена масла добавлена.\e[0m"
  end

  def add_reminder(car)
    print "Название: "
    title = gets.chomp
    print "Дата (ГГГГ-ММ-ДД): "
    date = gets.chomp
    print "Пробег (км): "
    mileage = gets.chomp.to_i
    print "Описание: "
    desc = gets.chomp
    print "Приоритет (низкий/средний/высокий): "
    priority = gets.chomp.downcase
    priority = "средний" unless ["низкий", "средний", "высокий"].include?(priority)
    car.add_reminder(title, date, mileage, desc, priority)
    save
    puts "\e[32m✅ Напоминание добавлено.\e[0m"
  end

  def show_history(car)
    puts "\n\e[36m📋 История для #{car}\e[0m"
    unless car.repairs.empty?
      puts "\e[33mРемонты:\e[0m"
      car.repairs.each do |r|
        puts "  #{r[:date]} | #{r[:mileage]} км | #{r[:description]} | #{r[:cost]} руб."
      end
    end
    unless car.oil_changes.empty?
      puts "\e[33mЗамены масла:\e[0m"
      car.oil_changes.each do |o|
        puts "  #{o[:date]} | #{o[:mileage]} км | #{o[:oil_type]} | #{o[:filter]} | #{o[:cost]} руб. | след. #{o[:next_mileage]} км"
      end
    end
    unless car.reminders.empty?
      puts "\e[33mНапоминания:\e[0m"
      car.reminders.each do |r|
        status = r[:done] ? "✅" : "⏳"
        puts "  #{status} #{r[:title]} | #{r[:date]} | #{r[:priority]} | #{r[:description]}"
      end
    end
  end

  def show_stats(car)
    total_repair = car.repairs.sum { |r| r[:cost] }
    total_oil = car.oil_changes.sum { |o| o[:cost] }
    total = total_repair + total_oil
    puts "\n\e[36m📊 Статистика для #{car}\e[0m"
    puts "  Всего ремонтов: #{car.repairs.size}"
    puts "  Всего замен масла: #{car.oil_changes.size}"
    puts "  Всего напоминаний: #{car.reminders.size}"
    puts "  Общая стоимость: #{total.round(2)} руб."
    puts "  Средняя стоимость ремонта: #{(total_repair / car.repairs.size).round(2)} руб." unless car.repairs.empty?
  end
end

def main
  manager = CarManager.new
  loop do
    puts "\n\e[36m🚗 Многомарочный дневник автомобиля (Ruby)\e[0m"
    puts "1. Выбрать/создать автомобиль"
    puts "2. Удалить автомобиль"
    puts "3. Показать все автомобили"
    puts "4. Выход"
    print "Выберите действие: "
    choice = gets.chomp
    case choice
    when "1"
      manager.list_cars
      if manager.cars.any?
        print "Введите номер автомобиля (или 0 для создания нового): "
        id_str = gets.chomp
        if id_str == "0"
          print "Марка: "
          brand = gets.chomp
          print "Модель: "
          model = gets.chomp
          print "Год: "
          year = gets.chomp.to_i
          print "VIN: "
          vin = gets.chomp
          print "Текущий пробег (км): "
          mileage = gets.chomp.to_i
          id = manager.add_car(brand, model, year, vin, mileage)
          puts "\e[32m✅ Автомобиль добавлен (ID: #{id})\e[0m"
          car = manager.get_car(id)
          manager.car_menu(car) if car
        elsif id_str.to_i > 0
          car = manager.get_car(id_str.to_i)
          if car
            manager.car_menu(car)
          else
            puts "\e[31m❌ Автомобиль не найден.\e[0m"
          end
        end
      else
        print "Марка: "
        brand = gets.chomp
        print "Модель: "
        model = gets.chomp
        print "Год: "
        year = gets.chomp.to_i
        print "VIN: "
        vin = gets.chomp
        print "Текущий пробег (км): "
        mileage = gets.chomp.to_i
        id = manager.add_car(brand, model, year, vin, mileage)
        puts "\e[32m✅ Автомобиль добавлен (ID: #{id})\e[0m"
        car = manager.get_car(id)
        manager.car_menu(car) if car
      end
    when "2"
      manager.list_cars
      print "Введите ID автомобиля для удаления: "
      id = gets.chomp.to_i
      if manager.delete_car(id)
        puts "\e[32m✅ Автомобиль удалён.\e[0m"
      else
        puts "\e[31m❌ Автомобиль не найден.\e[0m"
      end
    when "3"
      manager.list_cars
    when "4"
      puts "До свидания!"
      break
    else
      puts "\e[31mНеверный выбор.\e[0m"
    end
  end
end

main if __FILE__ == $0
