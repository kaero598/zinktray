(function () {
    const mailboxSelector = document.querySelector('.mailbox-selector');

    const messageList = document.querySelector('.message-list');
    const messageListEmptyMessage = messageList.querySelector('.message-list-empty');
    const messageListLoadingMessage = messageList.querySelector('.message-list-loading');
    const messageSelector = messageList.querySelector('.message-selector');

    const messageView = document.querySelector('main');

    var selectedMailbox;

    // App specific

    /**
     * @param {string} mailboxName
     */
    async function loadMailbox(mailboxName) {
        const mailboxes = await fetchMailboxList();

        var mailbox = mailboxes.find(
            (mailbox) => mailbox.name === mailboxName || (mailboxName === '' && mailbox.isAnonymous)
        );

        if (!mailbox) {
            mailbox = {
                id: mailboxName,
                name: mailboxName,
                isAnonymous: false,
            };
        }

        setSelectedMailbox(mailbox);

        await updateMessageList(mailbox.id);
    }

    /**
     * @param {string} mailboxId
     * @param {string} messageId
     */
    async function loadMessage(messageId) {
        setSelectedMessage(messageId);

        let message;

        try {
            message = await fetchMessageDetails(messageId);
        } catch (e) {
            resetMessageView();

            return
        }

        updateMessageView(message);
    }

    /**
     * @param {object} mailbox
     */
    function setSelectedMailbox(mailbox) {
        selectedMailbox = mailbox;

        mailboxSelector.textContent = mailbox.name;
    }

    /**
     * @param {string} messageId
     */
    function setSelectedMessage(messageId) {
        const activeMessageListItems = messageSelector.querySelectorAll('.active');

        if (activeMessageListItems.length) {
            Array.from(activeMessageListItems).forEach((node) => node.classList.remove('active'));
        }

        const lookupMessageListItem = messageSelector.querySelector(`[data-message-id="${messageId}"]`);

        if (lookupMessageListItem) {
            lookupMessageListItem.classList.add('active');
        }
    }

    /**
     * @param {string} mailboxId
     */
    async function updateMessageList(mailboxId) {
        messageListEmptyMessage.setAttribute('hidden', '');
        messageListLoadingMessage.removeAttribute('hidden');
        messageSelector.setAttribute('hidden', '');

        clearMessageList();

        messages = await fetchMessageList(mailboxId)

        if (!messages.length) {
            messageListEmptyMessage.removeAttribute('hidden');
            messageListLoadingMessage.setAttribute('hidden', '');
            messageSelector.setAttribute('hidden', '');

            return;
        }

        messages.forEach((message) => messageSelector.appendChild(createMessageListItem(message)));

        messageListEmptyMessage.setAttribute('hidden', '');
        messageListLoadingMessage.setAttribute('hidden', '');
        messageSelector.removeAttribute('hidden');
    }

    function clearMessageList() {
        Array.from(messageSelector.children).forEach((node) => messageSelector.removeChild(node));
    }

    /**
     * @param {object} message
     * @param {string} message.id
     * @param {string} message.subject
     * @param {string[]} message.to
     * @param {number} message.receivedAt
     * 
     * @returns Node
     */
    function createMessageListItem(message) {
        const subject = document.createElement('div');
        subject.classList.add('subject');

        if (message.subject === '') {
            subject.classList.add('empty-subject');
            subject.textContent = '<No subject>';
        } else {
            subject.textContent = message.subject;
        }

        const recipient = document.createElement('span');
        recipient.classList.add('recipient');
        recipient.textContent = 'To: ' + message.to.join(', ');

        const receiveTime = document.createElement('span');
        receiveTime.classList.add('receive-time');
        receiveTime.textContent = formatDate(new Date(message.receivedAt * 1000));

        const meta = document.createElement('div');
        meta.classList.add('meta');
        meta.appendChild(recipient);
        meta.appendChild(receiveTime);

        const encodedMailboxName = encodeURIComponent(selectedMailbox.name);
        const encodedMessageId = encodeURIComponent(message.id);

        const anchor = document.createElement('a');
        anchor.href = `#!/${encodedMailboxName}/${encodedMessageId}`;
        anchor.appendChild(subject);
        anchor.appendChild(meta);

        const item = document.createElement('li');
        item.setAttribute('data-message-id', message.id);

        item.appendChild(anchor);

        return item;
    }

    /**
     * @param {object} message
     */
    function updateMessageView(message) {
        clearMessageView();

        const subject = document.createElement('h2');
        subject.textContent = 'From: ' + message.subject;

        const from = document.createElement('p');
        from.textContent = 'From: ' + message.from.join(', ');

        const to = document.createElement('p');
        to.textContent = 'To: ' + message.to.join(', ');

        const date = document.createElement('p');
        date.textContent = 'Date: ' + formatDate(new Date(message.receivedAt * 1000));

        const header = document.createElement('section');

        header.appendChild(subject);
        header.appendChild(from);
        header.appendChild(to);
        header.appendChild(date);

        const contents = document.createElement('section');
        contents.classList.add('message-view-content');
        contents.textContent = message.rawBody;

        messageView.appendChild(header);
        messageView.appendChild(contents);
    }

    function resetMessageView() {
        clearMessageView();

        const message = document.createElement('p');
        message.textContent = 'Select message first.';

        messageView.appendChild(message);
    }

    function clearMessageView() {
        Array.from(messageView.children).forEach((node) => messageView.removeChild(node));
    }

    // Misc.

    /**
     * @param {Date} date 
     * 
     * @returns {string}
     */
    function formatDate(date) {
        const day = date.getDate();
        var month = date.getMonth();
        const year = date.getYear();

        if (month < 10) {
            month = '0' + month;
        }

        var result = `${day}.${month}`

        const currentYear = new Date().getYear();

        if (currentYear !== year) {
            const shortYear = year % 100;

            result += `.${shortYear}`;
        }

        return result;
    }

    // App flow

    /**
     * @param {string} hashUrl
     */
    function handleNavigation(hashUrl) {
        const path = hashUrl.substring(3);
        const segments = path.split('/', 3);

        if (segments.length > 2) {
            navigateTo();

            return;
        }

        [mailboxName, messageId] = segments;

        loadMailbox(mailboxName).then(() => {
            if (messageId !== undefined) {
                loadMessage(messageId);
            }
        });
    }

    /**
     * @param {string} [mailboxName]
     * @param {string} [messageId]
     */
    function navigateTo(mailboxName, messageId) {
        if (mailboxName === undefined) {
            location.hash = '#';

            return;
        }

        if (messageId === undefined) {
            location.hash = `#!/${mailboxName}`;

            return;
        }

        location.hash = `#!/${mailboxName}/${messageId}`;
    }

    // APIs

    /**
     * @returns {Promise<object[]>}
     */
    async function fetchMailboxList() {
        return fetch('/api/mailboxes/list').then((response) => response.json());
    }

    /**
     * @param {string} messageId
     * 
     * @returns {Promise<object>}
     */
    async function fetchMessageDetails(messageId) {
        const encodedMessageId = encodeURIComponent(messageId);

        return fetch(`/api/messages/details?message_id=${encodedMessageId}`).then((response) => response.json());
    }

    /**
     * @param {string} mailboxId
     * 
     * @returns {Promise<object[]>}
     */
    async function fetchMessageList(mailboxId) {
        const encodedMailboxId = encodeURIComponent(mailboxId);

        return fetch(`/api/messages/list?mailbox_id=${encodedMailboxId}`).then((response) => response.json());
    }

    // Event listeners

    document.addEventListener('DOMContentLoaded', function () {
        handleNavigation(location.hash);
    });

    window.addEventListener('hashchange', function (e) {
        handleNavigation(new URL(e.newURL).hash);
    });

    mailboxSelectorInput.addEventListener('keypress', function (e) {
        if (e.key !== 'Enter') {
            return;
        }

        navigateTo(mailboxSelectorInput.value.trim());
    });
})();