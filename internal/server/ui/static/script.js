const $ = sel => document.querySelector(sel);
const btn = $('#btn');
const idInput = $('#orderId');

btn.addEventListener('click', fetchOrder);
idInput.addEventListener('keydown', e => { if (e.key === 'Enter') fetchOrder(); });

async function fetchOrder() {
  const id = idInput.value.trim();
  if (!id) return;
  btn.disabled = true;
  try {
    const res = await fetch(`/orders/${encodeURIComponent(id)}`);
    if (!res.ok) {
      const msg = await res.text();
      alert(msg || 'Ошибка запроса');
      showNone();
      return;
    }
    const o = await res.json();
    fillOrder(o);
  } catch (e) {
    alert('Сеть недоступна или сервер не запущен');
    showNone();
  } finally {
    btn.disabled = false;
  }
}

function showNone() {
  $('#result').style.display = 'none';
  $('#itemsWrap').style.display = 'none';
  $('#raw').style.display = 'none';
}

function fillOrder(o) {
  $('#result').style.display = 'flex';
  $('#itemsWrap').style.display = 'block';
  $('#raw').style.display = 'block';

  $('#ord_uid').textContent = o.order_uid || '—';
  $('#ord_track').textContent = o.track_number || '—';
  $('#ord_entry').textContent = o.entry || '—';
  $('#ord_locale').textContent = o.locale || '—';
  $('#ord_customer').textContent = o.customer_id || '—';
  $('#ord_service').textContent = o.delivery_service || '—';
  $('#ord_date').textContent = o.date_created ? new Date(o.date_created).toLocaleString() : '—';

  const d = o.delivery || {};
  $('#del_name').textContent = d.name || '—';
  $('#del_phone').textContent = d.phone || '—';
  $('#del_addr').textContent = d.address || '—';
  $('#del_city').textContent = d.city || '—';
  $('#del_region').textContent = d.region || '—';
  $('#del_email').textContent = d.email || '—';
  $('#del_zip').textContent = d.zip || '—';

  const p = o.payment || {};
  $('#pay_tx').textContent = p.transaction || '—';
  $('#pay_amount').textContent = p.amount != null ? p.amount : '—';
  $('#pay_cur').textContent = p.currency || '—';
  $('#pay_bank').textContent = p.bank || '—';
  $('#pay_ship').textContent = p.delivery_cost != null ? p.delivery_cost : '—';
  $('#pay_goods').textContent = p.goods_total != null ? p.goods_total : '—';
  $('#pay_dt').textContent = p.payment_dt ? new Date(p.payment_dt).toLocaleString() : '—';

  const tbody = $('#itemsTable tbody');
  tbody.innerHTML = '';
  (o.items || []).forEach(it => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${it.chrt_id ?? ''}</td>
      <td>${escapeHtml(it.name ?? '')}</td>
      <td>${escapeHtml(it.brand ?? '')}</td>
      <td>${escapeHtml(it.size ?? '')}</td>
      <td>${it.price ?? ''}</td>
      <td>${it.total_price ?? ''}</td>
      <td>${it.status ?? ''}</td>`;
    tbody.appendChild(tr);
  });
  $('#itemsCount').textContent = `Всего: ${(o.items||[]).length}`;
  $('#rawJson').textContent = JSON.stringify(o, null, 2);
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, m => ({
    '&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'
  }[m]));
}