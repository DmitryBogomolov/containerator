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

    function getMode(row) {
        const sel = row.querySelector('.cmd-modes');
        return sel.options[sel.options.selectedIndex].value;
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
        const args = Object.assign({
            mode: getMode(row),
            tag: getTag(row)
        }, options);
        callManage(getRowId(row), args).then(
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
            const modelSel = row.querySelector('.cmd-modes');
            const tagsSel = row.querySelector('.cmd-tags');
            loadInfo(getRowId(row)).then(
                function ({ modes, tags }) {
                    addOptionItems(modelSel, modes || ['']);
                    addOptionItems(tagsSel, tags);
                },
                function (err) {
                    modelSel.classList.add('reduced-width');
                    addBadge(modelSel, err.message, true);
                    tagsSel.classList.add('reduced-width');
                    addBadge(tagsSel, err.message, true);
                }
            );
        });
    }

    function loadInfo(name) {
        return fetch('/api/info/' + name).then(function (response) {
            return response.ok
                ? response.json()
                : response.text().then(function (text) {
                    throw new Error(text);
                });
        });
    }

    function callManage(name, { mode, tag, force, remove }) {
        const formData = new FormData();
        if (mode) {
            formData.append('mode', mode);
        }
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
