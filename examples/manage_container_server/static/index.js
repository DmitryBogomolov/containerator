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
        const args = {
            tag: getTag(row),
            ...options,
        };
        const targetItemName = getRowId(row);
        commandManageContainer(targetItemName, args).then(
            (data) => {
                const content = `${data.name} ${data.image} ${data.tag}`;
                addBadge(button, content, false);
            },
            (err) => {
                addBadge(button, err.message, true);
            },
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
        const element = document.createElement('option');
        element.value = tag;
        element.textContent = tag;
        return element;
    }

    function addOptionItems(selectElement, optionValues) {
        for (const optionValue of optionValues) {
            const optionElement = createSelectOption(optionValue);
            selectElement.appendChild(optionElement);
        }
    }

    function attachHandlers() {
        const tableRows = document.querySelector('.table > tbody').children;
        for (const row of tableRows) {
            const targetItemName = getRowId(row);
            row.querySelector('.cmd-update').addEventListener('click', handleUpdateClick);
            row.querySelector('.cmd-remove').addEventListener('click', handleRemoveClick);
            const tagsSelect = row.querySelector('.cmd-tags');
            queryImageInfo(targetItemName).then(
                ({ tags }) => {
                    addOptionItems(tagsSelect, tags);
                },
                (err) => {
                    tagsSelect.classList.add('reduced-width');
                    addBadge(tagsSelect, err.message, true);
                },
            );
        }
    }

    function processJsonResponse(response) {
        return response.ok
            ? response.json()
            : response.text().then((text) => { throw new Error(text); });
    }

    function queryImageInfo(targetName) {
        return fetch(`/api/image-info/${targetName}`).then(processJsonResponse);
    }

    function commandManageContainer(targetName, payload) {
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        };
        return fetch(`/api/manage-container/${targetName}`, options).then(processJsonResponse);
    }

}());
