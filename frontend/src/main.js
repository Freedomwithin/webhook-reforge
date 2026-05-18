// WebhookReforge — frontend/src/main.js
// Wails bindings are exposed via window.go.main.App.*

let proxyRunning = false;
let currentMode = 'replay';

// ── Mode switching ──────────────────────────────────────────────

function setMode(mode) {
  currentMode = mode;

  document.getElementById('tab-replay').classList.toggle('active', mode === 'replay');
  document.getElementById('tab-proxy').classList.toggle('active', mode === 'proxy');

  const payloadSection = document.getElementById('payload-section');
  const proxyPanel = document.getElementById('proxy-panel');

  if (mode === 'replay') {
    payloadSection.style.display = 'flex';
    proxyPanel.classList.remove('visible');
  } else {
    payloadSection.style.display = 'none';
    proxyPanel.classList.add('visible');
  }
}

// ── Replay ──────────────────────────────────────────────────────

async function fireReplay() {
  const payload = document.getElementById('payload-editor').value.trim();
  const target  = document.getElementById('target').value.trim();
  const secret  = document.getElementById('secret').value.trim();
  const provider = document.getElementById('provider').value;

  if (!payload || !target || !secret) {
    addLogEntry('error', 'MISSING', 'Payload, Target URL, and Secret are all required.', provider);
    return;
  }

  // Validate JSON
  try {
    JSON.parse(payload);
    document.getElementById('payload-editor').classList.remove('error');
  } catch (e) {
    document.getElementById('payload-editor').classList.add('error');
    addLogEntry('error', 'JSON ERR', `Invalid JSON: ${e.message}`, provider);
    return;
  }

  // Write payload to a temp file via Go, then replay
  // We pass the raw JSON string; Go side handles temp file creation
  const btn = document.getElementById('fire-btn');
  btn.disabled = true;
  btn.textContent = 'Firing…';
  setStatus('active', 'firing');

  try {
    const result = await window.go.main.App.ReplayPayload(payload, target, secret, provider);
    const isError = result.toLowerCase().startsWith('error');
    addLogEntry(isError ? 'error' : 'success', isError ? 'FAIL' : 'OK', result, provider);
    setStatus(isError ? '' : 'active', isError ? 'error' : 'fired');
    setTimeout(() => setStatus('', 'idle'), 2000);
  } catch (e) {
    addLogEntry('error', 'ERR', String(e), provider);
    setStatus('', 'error');
  } finally {
    btn.disabled = false;
    btn.textContent = '⚡ Re-sign & Fire';
  }
}

// ── Proxy ───────────────────────────────────────────────────────

async function toggleProxy() {
  if (proxyRunning) {
    // No StopProxy yet — inform user
    addLogEntry('info', 'INFO', 'Proxy can be stopped by restarting the app. StopProxy coming in next release.', 'system');
    return;
  }

  const port   = parseInt(document.getElementById('proxy-port').value, 10);
  const target = document.getElementById('target').value.trim();
  const secret = document.getElementById('secret').value.trim();
  const provider = document.getElementById('provider').value;

  if (!target || !secret) {
    addLogEntry('error', 'MISSING', 'Target URL and Secret are required to start the proxy.', provider);
    return;
  }

  const btn = document.getElementById('proxy-btn');
  btn.disabled = true;
  setStatus('active', 'starting proxy');

  try {
    const result = await window.go.main.App.StartProxy(port, target, secret, provider);
    const isError = result.toLowerCase().startsWith('error');

    if (!isError) {
      proxyRunning = true;
      btn.textContent = 'Proxy Running ●';
      btn.style.background = 'var(--success)';
      document.getElementById('proxy-status').textContent = result;
      document.getElementById('proxy-status').classList.add('running');
      setStatus('active', `proxy :${port}`);
    } else {
      addLogEntry('error', 'FAIL', result, provider);
      setStatus('', 'error');
      btn.disabled = false;
    }

    addLogEntry(isError ? 'error' : 'info', isError ? 'FAIL' : 'PROXY', result, provider);
  } catch (e) {
    addLogEntry('error', 'ERR', String(e), provider);
    setStatus('', 'error');
    btn.disabled = false;
  }
}

// ── Payload helpers ─────────────────────────────────────────────

function formatPayload() {
  const editor = document.getElementById('payload-editor');
  try {
    const parsed = JSON.parse(editor.value);
    editor.value = JSON.stringify(parsed, null, 2);
    editor.classList.remove('error');
  } catch (e) {
    editor.classList.add('error');
  }
}

// ── Event Log ───────────────────────────────────────────────────

function addLogEntry(type, badge, message, provider) {
  const log = document.getElementById('event-log');

  // Remove empty state
  const empty = document.getElementById('log-empty');
  if (empty) empty.remove();

  const now = new Date();
  const time = now.toLocaleTimeString('en-US', { hour12: false });

  const entry = document.createElement('div');
  entry.className = `log-entry ${type}`;
  entry.innerHTML = `
    <div class="log-entry-header">
      <span class="log-badge">${badge}</span>
      <span class="log-time">${time}</span>
      <span class="log-provider">${provider || ''}</span>
    </div>
    <div class="log-message">${escapeHtml(message)}</div>
  `;

  // Prepend so newest is at top
  log.insertBefore(entry, log.firstChild);
}

function clearLog() {
  const log = document.getElementById('event-log');
  log.innerHTML = '<div class="log-empty" id="log-empty">No events yet. Fire a webhook to get started.</div>';
  setStatus('', 'idle');
}

// ── Status indicator ────────────────────────────────────────────

function setStatus(dotClass, label) {
  const dot = document.getElementById('status-dot');
  const lbl = document.getElementById('status-label');
  dot.className = 'status-dot' + (dotClass ? ` ${dotClass}` : '');
  lbl.textContent = label;
}

// ── Utilities ───────────────────────────────────────────────────

function escapeHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
