document.addEventListener('DOMContentLoaded', function() {
    const API_BASE = '/';

    // Load items on page load
    loadItems();

    // Add item form
    const addForm = document.getElementById('add-form');
    addForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const formData = {
            type: document.getElementById('type').value,
            amount: parseInt(document.getElementById('amount').value),
            date: document.getElementById('date').value,
            category: document.getElementById('category').value
        };

        try {
            const response = await fetch(API_BASE + 'items', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            });

            if (response.ok) {
                addForm.reset();
                loadItems(); // Reload table
                alert('Запись добавлена!');
            } else {
                throw new Error('Ошибка при добавлении');
            }
        } catch (error) {
            alert('Ошибка: ' + error.message);
        }
    });

    // Apply filters and sorting
    const applyFiltersBtn = document.getElementById('apply-filters');
    applyFiltersBtn.addEventListener('click', function() {
        loadItems();
    });

    // Analytics form
    const analyticsForm = document.getElementById('analytics-form');
    analyticsForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const from = document.getElementById('analytics-from').value;
        const to = document.getElementById('analytics-to').value;
        await loadAnalytics(from, to);
    });

    // Export buttons
    document.getElementById('export-items-csv').addEventListener('click', exportItemsCSV);
    document.getElementById('export-analytics-csv').addEventListener('click', exportAnalyticsCSV);

    async function loadItems() {
        const sortBy = document.getElementById('sort-by').value;
        const params = new URLSearchParams();
        if (sortBy) params.append('sort_by', sortBy);

        const response = await fetch(API_BASE + 'items?' + params.toString());
        if (!response.ok) throw new Error('Ошибка загрузки записей');

        const items = await response.json();
        displayItems(items);
        applyClientSideFilters(items);
    }

    function displayItems(items) {
        const tbody = document.querySelector('#items-table tbody');
        tbody.innerHTML = '';

        items.forEach(item => {
            const row = tbody.insertRow();
            row.innerHTML = `
                <td>${item.id}</td>
                <td>${item.type}</td>
                <td>${item.amount}</td>
                <td>${item.date}</td>
                <td>${item.category}</td>
                <td>
                    <button onclick="editItem(${item.id}, '${item.type}', ${item.amount}, '${item.date}', '${item.category}')">Редактировать</button>
                    <button onclick="deleteItem(${item.id})">Удалить</button>
                </td>
            `;
        });
    }

    function applyClientSideFilters(allItems) {
        const filterType = document.getElementById('filter-type').value;
        const filterCategory = document.getElementById('filter-category').value.toLowerCase();
        const filterDateFrom = document.getElementById('filter-date-from').value;
        const filterDateTo = document.getElementById('filter-date-to').value;

        let filteredItems = allItems;

        if (filterType) {
            filteredItems = filteredItems.filter(item => item.type === filterType);
        }

        if (filterCategory) {
            filteredItems = filteredItems.filter(item => item.category.toLowerCase().includes(filterCategory));
        }

        if (filterDateFrom) {
            filteredItems = filteredItems.filter(item => item.date >= filterDateFrom);
        }

        if (filterDateTo) {
            filteredItems = filteredItems.filter(item => item.date <= filterDateTo);
        }

        displayItems(filteredItems);
    }

    // Edit item (simple prompt for now, can be improved with modal)
    window.editItem = async function(id, type, amount, date, category) {
        const newType = prompt('Новый тип:', type);
        const newAmount = prompt('Новая сумма:', amount);
        const newDate = prompt('Новая дата (YYYY-MM-DD):', date);
        const newCategory = prompt('Новая категория:', category);

        if (newType && newAmount && newDate && newCategory) {
            const updateData = {
                type: newType,
                amount: parseInt(newAmount),
                date: newDate,
                category: newCategory
            };

            const response = await fetch(API_BASE + `items/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(updateData)
            });

            if (response.ok) {
                loadItems();
                alert('Запись обновлена!');
            } else {
                alert('Ошибка обновления');
            }
        }
    };

    // Delete item
    window.deleteItem = async function(id) {
        if (confirm('Удалить запись?')) {
            const response = await fetch(API_BASE + `items/${id}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                loadItems();
                alert('Запись удалена!');
            } else {
                alert('Ошибка удаления');
            }
        }
    };

    async function loadAnalytics(from, to) {
        const params = new URLSearchParams();
        if (from) params.append('from', from);
        if (to) params.append('to', to);

        const response = await fetch(API_BASE + 'analytics?' + params.toString());
        if (!response.ok) throw new Error('Ошибка загрузки аналитики');

        const analytics = await response.json();
        displayAnalytics(analytics);
        drawChart(analytics);
    }

    function displayAnalytics(analytics) {
        const tbody = document.querySelector('#analytics-table tbody');
        tbody.innerHTML = '';

        if (analytics.length === 0) {
            tbody.innerHTML = '<tr><td colspan="2">Нет данных</td></tr>';
            return;
        }

        const agg = analytics[0].aggregated_data;
        const rows = [
            ['Сумма', agg.sum],
            ['Среднее', agg.average],
            ['Количество', agg.count],
            ['Медиана', agg.median],
            ['90-й перцентиль', agg.percentile_90],
        ];

        rows.forEach(([label, value]) => {
            const row = tbody.insertRow();
            row.innerHTML = `<td>${label}</td><td>${value}</td>`;
        });
    }

    function drawChart(analytics) {
        const canvas = document.getElementById('analytics-chart');
        const ctx = canvas.getContext('2d');

        // Simple bar chart for aggregated values
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        if (analytics.length === 0) return;

        // Set canvas size to fit 5 bars
        canvas.width = 400;
        canvas.height = 200;

        const agg = analytics[0].aggregated_data;
        // Bars for sum, average (scaled), median (scaled), count (scaled), percentile_90 (scaled)
        const scaleFactor = 10; // Scale smaller values for visibility
        const data = [
            agg.sum,
            agg.average * scaleFactor,
            agg.median * scaleFactor,
            agg.count * scaleFactor,
            agg.percentile_90 * scaleFactor
        ];
        const labels = ['Сумма', 'Среднее', 'Медиана', 'Количество', '90-й перцентиль'];
        const barWidth = 40;
        const spacing = 60;
        const maxHeight = 150;
        const maxValue = Math.max(...data);

        data.forEach((value, index) => {
            const height = (value / maxValue) * maxHeight;
            const x = index * spacing + 20;
            const y = canvas.height - height - 20;

            ctx.fillStyle = '#007bff';
            ctx.fillRect(x, y, barWidth, height);

            ctx.fillStyle = '#333';
            ctx.font = '12px Arial';
            ctx.fillText(labels[index], x, canvas.height - 5);
            // Display actual values (unscaled where appropriate)
            let actualValue;
            if (index === 0) {
                actualValue = value.toFixed(0);
            } else if (index === 3) {
                actualValue = (value / scaleFactor).toFixed(0);
            } else {
                actualValue = (value / scaleFactor).toFixed(2);
            }
            ctx.fillText(actualValue, x, y - 5);
        });
    }

    async function exportItemsCSV() {
        const sortBy = document.getElementById('sort-by').value;
        const params = new URLSearchParams();
        if (sortBy) params.append('sort_by', sortBy);

        const response = await fetch(API_BASE + 'items/csv?' + params.toString());
        if (response.ok) {
            const blob = await response.blob();
            downloadBlob(blob, 'items.csv');
        } else {
            alert('Ошибка экспорта');
        }
    }

    async function exportAnalyticsCSV() {
        const from = document.getElementById('analytics-from').value;
        const to = document.getElementById('analytics-to').value;
        const params = new URLSearchParams();
        if (from) params.append('from', from);
        if (to) params.append('to', to);

        const response = await fetch(API_BASE + 'analytics/csv?' + params.toString());
        if (response.ok) {
            const blob = await response.blob();
            downloadBlob(blob, 'analytics.csv');
        } else {
            alert('Ошибка экспорта');
        }
    }

    function downloadBlob(blob, filename) {
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }
});
