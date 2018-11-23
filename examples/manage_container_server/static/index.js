(function () {
    const table = document.querySelector('#projects');
    const body = table.tBodies[0];

    function refreshItems() {
        return fetch('/api/projects')
            .then(response => response.json())
            .then((items) => {
                items.forEach((item) => {
                    const row = document.createElement('tr');
                    row.appendChild(createNameColumn(item.name));
                    row.appendChild(createUpdateColumn(item.id));
                    row.appendChild(createForceColumn(item.id));
                    row.appendChild(createRemoveColumn(item.id));
                    body.appendChild(row);
                });
            })
            .catch((err) => {
                console.error(err);
            });
    }

    function createNameColumn(name) {
        const column = document.createElement('td');
        column.textContent = name;
        return column;
    }

    function createUpdateColumn(id) {
        const column = document.createElement('td');
        const button = document.createElement('button');
        button.textContent = 'Do it';
        button.addEventListener('click', () => {
        });
        column.appendChild(button);
        return column;
    }

    function createForceColumn(id) {
        const column = document.createElement('td');
        const input = document.createElement('input');
        input.setAttribute('type', 'checkbox');
        column.appendChild(input);
        return column;
    }

    function createRemoveColumn(id) {
        const column = document.createElement('td');
        const button = document.createElement('button');
        button.textContent = 'Do it';
        button.addEventListener('click', () => {
        });
        column.appendChild(button);
        return column;
    }

    refreshItems();

}());
