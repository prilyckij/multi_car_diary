<?php
// multi_car_diary.php — PHP версия

$dataFile = 'cars.json';

function loadCars() {
    global $dataFile;
    if (file_exists($dataFile)) {
        $json = file_get_contents($dataFile);
        return json_decode($json, true) ?: [];
    }
    return [];
}

function saveCars($cars) {
    global $dataFile;
    file_put_contents($dataFile, json_encode($cars, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
}

$cars = loadCars();

function color($text, $code) {
    return "\033[{$code}m{$text}\033[0m";
}

function listCars($cars) {
    if (empty($cars)) {
        echo color("Нет автомобилей.\n", '33');
        return;
    }
    printf(color("%-4s %-15s %-15s %-6s %-10s\n", '36'), "ID", "Марка", "Модель", "Год", "Пробег");
    echo str_repeat("-", 55) . "\n";
    foreach ($cars as $c) {
        printf("%-4d %-15s %-15s %-6d %-10d\n", $c['id'], $c['brand'], $c['model'], $c['year'], $c['mileage']);
    }
}

function getCar($cars, $id) {
    foreach ($cars as &$c) {
        if ($c['id'] == $id) return $c;
    }
    return null;
}

function addRepair(&$car) {
    echo "Дата (ГГГГ-ММ-ДД): ";
    $date = trim(fgets(STDIN));
    echo "Пробег (км): ";
    $mileage = (int) trim(fgets(STDIN));
    echo "Описание: ";
    $desc = trim(fgets(STDIN));
    echo "Стоимость (руб): ";
    $cost = (float) trim(fgets(STDIN));
    $car['repairs'][] = ['date' => $date, 'mileage' => $mileage, 'description' => $desc, 'cost' => $cost];
    if ($mileage > $car['mileage']) $car['mileage'] = $mileage;
    echo color("✅ Ремонт добавлен.\n", '32');
}

function addOilChange(&$car) {
    echo "Дата (ГГГГ-ММ-ДД): ";
    $date = trim(fgets(STDIN));
    echo "Пробег (км): ";
    $mileage = (int) trim(fgets(STDIN));
    echo "Тип масла: ";
    $oil = trim(fgets(STDIN));
    echo "Фильтр (артикул): ";
    $filter = trim(fgets(STDIN));
    echo "Стоимость (руб): ";
    $cost = (float) trim(fgets(STDIN));
    echo "Интервал (км, по умолч. 10000): ";
    $intervalStr = trim(fgets(STDIN));
    $interval = $intervalStr ? (int)$intervalStr : 10000;
    $car['oil_changes'][] = [
        'date' => $date, 'mileage' => $mileage, 'oil_type' => $oil,
        'filter' => $filter, 'cost' => $cost, 'interval' => $interval,
        'next_mileage' => $mileage + $interval
    ];
    if ($mileage > $car['mileage']) $car['mileage'] = $mileage;
    echo color("✅ Замена масла добавлена.\n", '32');
}

function addReminder(&$car) {
    echo "Название: ";
    $title = trim(fgets(STDIN));
    echo "Дата (ГГГГ-ММ-ДД): ";
    $date = trim(fgets(STDIN));
    echo "Пробег (км): ";
    $mileage = (int) trim(fgets(STDIN));
    echo "Описание: ";
    $desc = trim(fgets(STDIN));
    echo "Приоритет (низкий/средний/высокий): ";
    $priority = strtolower(trim(fgets(STDIN)));
    if (!in_array($priority, ['низкий', 'средний', 'высокий'])) $priority = 'средний';
    $car['reminders'][] = [
        'title' => $title, 'date' => $date, 'mileage' => $mileage,
        'description' => $desc, 'priority' => $priority, 'done' => false
    ];
    echo color("✅ Напоминание добавлено.\n", '32');
}

function showHistory($car) {
    echo "\n" . color("📋 История для {$car['brand']} {$car['model']} {$car['year']}\n", '36');
    if (!empty($car['repairs'])) {
        echo color("Ремонты:\n", '33');
        foreach ($car['repairs'] as $r) {
            echo "  {$r['date']} | {$r['mileage']} км | {$r['description']} | {$r['cost']} руб.\n";
        }
    }
    if (!empty($car['oil_changes'])) {
        echo color("Замены масла:\n", '33');
        foreach ($car['oil_changes'] as $o) {
            echo "  {$o['date']} | {$o['mileage']} км | {$o['oil_type']} | {$o['filter']} | {$o['cost']} руб. | след. {$o['next_mileage']} км\n";
        }
    }
    if (!empty($car['reminders'])) {
        echo color("Напоминания:\n", '33');
        foreach ($car['reminders'] as $r) {
            $status = $r['done'] ? "✅" : "⏳";
            echo "  {$status} {$r['title']} | {$r['date']} | {$r['priority']} | {$r['description']}\n";
        }
    }
}

function showStats($car) {
    $totalRepair = array_sum(array_column($car['repairs'], 'cost'));
    $totalOil = array_sum(array_column($car['oil_changes'], 'cost'));
    $total = $totalRepair + $totalOil;
    echo "\n" . color("📊 Статистика для {$car['brand']} {$car['model']} {$car['year']}\n", '36');
    echo "  Всего ремонтов: " . count($car['repairs']) . "\n";
    echo "  Всего замен масла: " . count($car['oil_changes']) . "\n";
    echo "  Всего напоминаний: " . count($car['reminders']) . "\n";
    echo "  Общая стоимость: " . number_format($total, 2) . " руб.\n";
    if (!empty($car['repairs'])) {
        echo "  Средняя стоимость ремонта: " . number_format($totalRepair / count($car['repairs']), 2) . " руб.\n";
    }
}

function carMenu(&$cars, $carId) {
    global $cars;
    $car = &$cars[array_search($carId, array_column($cars, 'id'))];
    while (true) {
        echo "\n" . color("🚗 {$car['brand']} {$car['model']} {$car['year']} (Пробег: {$car['mileage']} км) — меню автомобиля\n", '36');
        echo "1. Добавить ремонт\n";
        echo "2. Добавить замену масла\n";
        echo "3. Добавить напоминание\n";
        echo "4. Показать историю\n";
        echo "5. Показать статистику\n";
        echo "6. Назад\n";
        echo "Выберите действие: ";
        $choice = trim(fgets(STDIN));
        switch ($choice) {
            case "1": addRepair($car); break;
            case "2": addOilChange($car); break;
            case "3": addReminder($car); break;
            case "4": showHistory($car); break;
            case "5": showStats($car); break;
            case "6": return;
            default: echo color("Неверный выбор.\n", '31');
        }
        saveCars($cars);
    }
}

function createNewCar(&$cars) {
    echo "Марка: ";
    $brand = trim(fgets(STDIN));
    echo "Модель: ";
    $model = trim(fgets(STDIN));
    echo "Год: ";
    $year = (int) trim(fgets(STDIN));
    echo "VIN: ";
    $vin = trim(fgets(STDIN));
    echo "Текущий пробег (км): ";
    $mileage = (int) trim(fgets(STDIN));
    $id = count($cars) + 1;
    $cars[] = [
        'id' => $id, 'brand' => $brand, 'model' => $model, 'year' => $year,
        'vin' => $vin, 'mileage' => $mileage, 'repairs' => [], 'oil_changes' => [], 'reminders' => []
    ];
    saveCars($cars);
    echo color("✅ Автомобиль добавлен (ID: $id)\n", '32');
    carMenu($cars, $id);
}

function main() {
    global $cars;
    while (true) {
        echo "\n" . color("🚗 Многомарочный дневник автомобиля (PHP)\n", '36');
        echo "1. Выбрать/создать автомобиль\n";
        echo "2. Удалить автомобиль\n";
        echo "3. Показать все автомобили\n";
        echo "4. Выход\n";
        echo "Выберите действие: ";
        $choice = trim(fgets(STDIN));
        switch ($choice) {
            case "1":
                listCars($cars);
                if (!empty($cars)) {
                    echo "Введите номер автомобиля (или 0 для создания нового): ";
                    $idStr = trim(fgets(STDIN));
                    if ($idStr == "0") {
                        createNewCar($cars);
                    } elseif (is_numeric($idStr)) {
                        $carId = (int) $idStr;
                        if (in_array($carId, array_column($cars, 'id'))) {
                            carMenu($cars, $carId);
                        } else {
                            echo color("❌ Автомобиль не найден.\n", '31');
                        }
                    }
                } else {
                    createNewCar($cars);
                }
                break;
            case "2":
                listCars($cars);
                echo "Введите ID автомобиля для удаления: ";
                $id = (int) trim(fgets(STDIN));
                $index = array_search($id, array_column($cars, 'id'));
                if ($index !== false) {
                    array_splice($cars, $index, 1);
                    saveCars($cars);
                    echo color("✅ Автомобиль удалён.\n", '32');
                } else {
                    echo color("❌ Автомобиль не найден.\n", '31');
                }
                break;
            case "3":
                listCars($cars);
                break;
            case "4":
                echo "До свидания!\n";
                exit(0);
            default:
                echo color("Неверный выбор.\n", '31');
        }
    }
}

main();
?>
