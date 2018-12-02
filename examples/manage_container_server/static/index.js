(function () {
    attachHandlers();

    function getRowId(row) {
        return row.dataset.id;
    }

    function findRow(node) {
        for (let x = node; x; x = x.parentNode) {
            if (getRowId(x)) {
                return x;
            }
        }
    }

    function getTag(row) {
        const sel = row.querySelector('.cmd-tags');
        return sel.options[sel.options.selectedIndex].value;
    }

    function handleUpdateClick() {
        const row = findRow(this);
        callManage(getRowId(row), {
            tag: getTag(row),
            force: row.querySelector('.cmd-force').checked
        });
    }

    function handleRemoveClick() {
        const row = findRow(this);
        callManage(getRowId(row), {
            tag: getTag(row),
            remove: true
        });
    }

    function createSelectOption(tag) {
        const opt = document.createElement('option');
        opt.value = tag;
        opt.textContent = tag;
        return opt;
    }

    function attachHandlers() {
        const body = document.querySelector('.table > tbody');
        Array.from(body.children).forEach(function (row) {
            row.querySelector('.cmd-update').addEventListener('click', handleUpdateClick);
            row.querySelector('.cmd-remove').addEventListener('click', handleRemoveClick);
            loadTags(getRowId(row)).then(function (tags) {
                const sel = row.querySelector('.cmd-tags');
                tags.map(createSelectOption).forEach(function (item) {
                    sel.appendChild(item);
                });
            });
        });
    }

    function loadTags(name) {
        return fetch('/api/tags/' + name)
            .then(function (response) {
                return response.json();
            })
            .catch(function (err) {
                console.error(err);
                return [];
            });
    }

    function callManage(name, { tag, force, remove }) {
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
        fetch('/api/manage/' + name, options)
            .then(function (response) {
                return response.json();
            })
            .catch(function (err) {
                console.error(err);
            });
    }

}());
