document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('shorten-form');
    const fullUrlInput = document.getElementById('full-url');
    const resultContainer = document.getElementById('result-container');
    const shortUrlLink = document.getElementById('short-url-link');
    const copyBtn = document.getElementById('copy-btn');
    const urlsList = document.getElementById('urls-list');

    // Initial fetch
    fetchURLs();

    // Check for errors in URL params
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('error') === 'not_found') {
        showToast('ERROR: LINK NOT FOUND');
        // Clean up the URL without refreshing
        window.history.replaceState({}, document.title, "/");
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const fullUrl = fullUrlInput.value;
        const btn = form.querySelector('button');
        const originalBtnText = btn.innerHTML;

        try {
            btn.disabled = true;
            btn.innerText = 'WAIT...';

            const response = await fetch('/api/url', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ full_url: fullUrl })
            });

            if (!response.ok) throw new Error('Failed to shorten URL');

            const data = await response.json();
            const shortUrl = `${window.location.origin}/${data.code}`;
            
            shortUrlLink.href = shortUrl;
            shortUrlLink.innerText = shortUrl;
            resultContainer.classList.remove('hidden');
            
            fullUrlInput.value = '';
            fetchURLs(); // Refresh list
        } catch (err) {
            alert(err.message);
        } finally {
            btn.disabled = false;
            btn.innerHTML = originalBtnText;
        }
    });

    copyBtn.addEventListener('click', () => {
        const text = shortUrlLink.innerText;
        navigator.clipboard.writeText(text).then(() => {
            showToast('LINK COPIED TO CLIPBOARD');
            const originalSvg = copyBtn.innerHTML;
            copyBtn.innerHTML = '<span style="font-size: 0.7rem; color: #0070ff">COPIED!</span>';
            setTimeout(() => {
                copyBtn.innerHTML = originalSvg;
            }, 2000);
        });
    });

    function showToast(message) {
        const toast = document.getElementById('toast');
        toast.innerText = message;
        toast.classList.add('show');
        setTimeout(() => {
            toast.classList.remove('show');
        }, 3000);
    }

    async function fetchURLs() {
        try {
            const response = await fetch('/api/url');
            if (!response.ok) throw new Error('Failed to fetch URLs');
            
            const urls = await response.json();
            renderURLs(urls);
        } catch (err) {
            urlsList.innerHTML = `<div class="loading-state">Error loading links: ${err.message}</div>`;
        }
    }

    function renderURLs(urls) {
        if (!urls || urls.length === 0) {
            urlsList.innerHTML = '<div class="loading-state">No links shortened yet. Start above!</div>';
            return;
        }

        urlsList.innerHTML = urls.reverse().map(url => `
            <div class="url-card">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <a href="${window.location.origin}/${url.code}" target="_blank" class="card-code">/${url.code}</a>
                    <div class="card-clicks">${url.click_count || 0} CLICKS</div>
                </div>
                <a href="${url.full_url}" target="_blank" class="card-full" title="${url.full_url}">${url.full_url}</a>
                <div class="card-actions">
                    <button class="icon-btn" onclick="copyToClipboard('${window.location.origin}/${url.code}')" title="Copy">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                    </button>
                    <button class="icon-btn" onclick="deleteURL('${url.code}')" title="Delete">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                    </button>
                </div>
            </div>
        `).join('');
    }
});

// Global functions for inline event handlers
async function copyToClipboard(text) {
    await navigator.clipboard.writeText(text);
    // Use the toast function defined in DOMContentLoaded context
    const toast = document.getElementById('toast');
    toast.innerText = 'LINK COPIED TO CLIPBOARD';
    toast.classList.add('show');
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

async function deleteURL(code) {
    if (!confirm('Are you sure you want to delete this link?')) return;
    
    try {
        const response = await fetch(`/api/url/${code}`, { method: 'DELETE' });
        if (response.ok) {
            location.reload(); // Simple refresh for now
        }
    } catch (err) {
        alert('Delete failed');
    }
}
