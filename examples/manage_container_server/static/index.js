(function () {
    const table = document.querySelector('#projects');
    const body = table.tBodies[0];

    function refreshItems() {
        fetch('/api/projects')
            .then(function (response) {
                return response.json();
            })
            .then(function (names) {
                body.innerHTML = '';
                names.forEach(function (name) {
                    const row = document.createElement('tr');
                    const ctx = {};
                    row.appendChild(createNameColumn(name));
                    row.appendChild(createUpdateColumn(name, ctx));
                    row.appendChild(createForceColumn(name, ctx));
                    row.appendChild(createRemoveColumn(name));
                    body.appendChild(row);
                });
            })
            .catch(function (err) {
                console.error(err);
            });
    }

    function findRow(node) {
        for (let x = node; x; x = x.parentNode) {
            if (x.dataset.id) {
                return x;
            }
        }
    }

    function handleUpdateClick() {
        const row = findRow(this);
        const name = row.dataset.id;
        callManage(name, {
            force: row.querySelector('.cmd-force').checked
        });
    }

    function handleRemoveClick() {
        const row = findRow(this);
        const name = row.dataset.id;
        callManage(name, { remove: true });
    }

    function attachHandlers() {
        const body = document.querySelector('#projects > tbody');
        Array.from(body.children).forEach(function (row) {
            row.querySelector('.cmd-update').addEventListener('click', handleUpdateClick);
            row.querySelector('.cmd-remove').addEventListener('click', handleRemoveClick);
        });
    }

    attachHandlers();

    function callManage(name, { tag, force, remove }) {
        console.log('MANAGE', name, force, remove);
        return;

        const formData = new FormData();
        if (tag) {
            formData.append('tag', tag);
        }
        if (force) {
            formData.append('force', true);
        }
        if (remove) {
            formData.append('remove', true);
        }
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: formData
        };
        fetch('/manage/' + name, options)
            .then(function (response) {
                return response.json();
            })
            .catch(function (err) {
                console.error(err);
            });
    }

    function createNameColumn(name) {
        const column = document.createElement('td');
        column.textContent = name;
        return column;
    }

    function createUpdateColumn(name, ctx) {
        const column = document.createElement('td');
        const button = document.createElement('button');
        button.textContent = 'Do it';
        button.addEventListener('click', function () {
            callManage(name, { force: ctx.isForced() });
        });
        column.appendChild(button);
        return column;
    }

    function createForceColumn(name, ctx) {
        const column = document.createElement('td');
        const input = document.createElement('input');
        input.setAttribute('type', 'checkbox');
        column.appendChild(input);
        ctx.isForced = function () {
            return input.checked;
        };
        return column;
    }

    function createRemoveColumn(name) {
        const column = document.createElement('td');
        const button = document.createElement('button');
        button.textContent = 'Do it';
        button.addEventListener('click', function () {
            callManage(name, { remove: true });
        });
        column.appendChild(button);
        return column;
    }

    // refreshItems();

}());
