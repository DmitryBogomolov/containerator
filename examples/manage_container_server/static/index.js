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

    function addBadge(node, content, isError) {
        const parent = node.parentNode;
        let badge = parent.querySelector('.badge');
        if (badge) {
            parent.removeChild(badge);
        }
        badge = document.createElement('span');
        badge.setAttribute('class', 'badge badge-' + (isError ?  'danger' : 'success'));
        badge.textContent = isError ? 'Error' : 'Success';
        $(badge).popover({
            trigger: 'hover',
            title: isError ? 'Error' : 'Success',
            content: content
        });
        parent.appendChild(badge);
    }

    function handleManageClick(button, options) {
        const row = findRow(button);
        callManage(getRowId(row), Object.assign({ tag: getTag(row) }, options)).then(
            function (data) {
                const content = `${data.name} ${data.image} ${data.tag}`;
                addBadge(button, content, false);
            },
            function (err) {
                addBadge(button, err.message, true);
            }
        );
    }

    function handleUpdateClick() {
        handleManageClick(this, {
            force: findRow(this).querySelector('.cmd-force').checked
        });
    }

    function handleRemoveClick() {
        handleManageClick(this, {
            remove: true
        })
    }

    function createSelectOption(tag) {
        const opt = document.createElement('option');
        opt.value = tag;
        opt.textContent = tag;
        return opt;
    }

    function addOptionItems(sel, tags) {
        tags.map(createSelectOption).forEach(function (item) {
            sel.appendChild(item);
        });
    }

    function attachHandlers() {
        const body = document.querySelector('.table > tbody');
        Array.from(body.children).forEach(function (row) {
            row.querySelector('.cmd-update').addEventListener('click', handleUpdateClick);
            row.querySelector('.cmd-remove').addEventListener('click', handleRemoveClick);
            const sel = row.querySelector('.cmd-tags');
            loadTags(getRowId(row)).then(
                function (tags) {
                    addOptionItems(sel, tags);
                },
                function (err) {
                    sel.classList.add('reduced-width');
                    addBadge(sel, err.message, true);
                }
            );
        });
    }

    function loadTags(name) {
        return fetch('/api/tags/' + name).then(function (response) {
            return response.ok
                ? response.json()
                : response.text().then(function (text) {
                    throw new Error(text);
                });
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
        return fetch('/api/manage/' + name, options).then(function (response) {
            return response.ok
                ? response.json()
                : response.text().then(function (text) {
                    throw new Error(text);
                });
        });
    }

}());
