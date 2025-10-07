(() => {
  const uploadForm = document.getElementById('uploadForm');
  const imageFile = document.getElementById('imageFile');
  const dropzone = document.getElementById('dropzone');
  const taskSelect = document.getElementById('task');
  const widthInput = document.getElementById('width');
  const heightInput = document.getElementById('height');
  const watermarkInput = document.getElementById('watermark');
  const uploadMsg = document.getElementById('uploadMsg');
  const imagesList = document.getElementById('images');
  const tmpl = document.getElementById('imageItemTmpl');
  const toasts = document.getElementById('toasts');

  const API = {
    upload: '/api/upload',
    info: (id) => `/api/image/info/${id}`,
    image: (id) => `/api/image/${id}`,
    delete: (id) => `/api/image/${id}`
  };

  const statusMap = {
    in_progress: ['in progress', false],
    finished: ['finished', true]
  };

  // Toasts
  function showToast(text, timeout = 2500){
    const el = document.createElement('div');
    el.className = 'toast';
    el.textContent = text;
    toasts.appendChild(el);
    setTimeout(()=>{ el.remove(); }, timeout);
  }

  function setTaskFieldsVisibility() {
    const task = taskSelect.value;
    document.querySelectorAll('[data-task]')
      .forEach(el => {
        const requiredTask = el.getAttribute('data-task');
        if (requiredTask === 'resize') {
          el.classList.toggle('hidden', task !== 'resize');
        } else if (requiredTask === 'watermark') {
          el.classList.toggle('hidden', task !== 'watermark');
        }
      });
  }

  taskSelect.addEventListener('change', setTaskFieldsVisibility);
  setTaskFieldsVisibility();

  uploadForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    uploadMsg.textContent = '';
    if (!imageFile.files[0]) return;

    const file = imageFile.files[0];
    const contentType = file.type;

    const task = taskSelect.value;
    const metadata = {
      task,
      content_type: contentType,
      watermark_string: watermarkInput.value || '',
      resize: {
        width: Number(widthInput.value || 0),
        height: Number(heightInput.value || 0)
      }
    };

    const fd = new FormData();
    fd.append('image', file);
    fd.append('metadata', JSON.stringify(metadata));

    try {
      const resp = await fetch(API.upload, { method: 'POST', body: fd });
      const json = await resp.json();
      if (!resp.ok) throw new Error(json.message || 'Upload failed');
      const id = json.result;
      uploadMsg.textContent = 'Задача отправлена. ID: ' + id;
      addImageItem(id);
      showToast('Файл загружен, начата обработка');
    } catch (err) {
      uploadMsg.textContent = 'Ошибка: ' + err.message;
      showToast('Ошибка: ' + err.message, 3500);
    }
  });

  // Drag & Drop
  ['dragenter','dragover'].forEach(evt => dropzone.addEventListener(evt, (e)=>{ e.preventDefault(); dropzone.classList.add('active'); }))
  ;['dragleave','drop'].forEach(evt => dropzone.addEventListener(evt, (e)=>{ e.preventDefault(); dropzone.classList.remove('active'); }))
  dropzone.addEventListener('drop', (e)=>{
    const files = Array.from(e.dataTransfer.files || []);
    if (!files.length) return;
    const img = files.find(f => /image\//.test(f.type));
    if (img) {
      imageFile.files = new DataTransfer([img]).files;
      showToast('Файл выбран: ' + img.name);
    }
  });

  function addImageItem(id) {
    const node = tmpl.content.firstElementChild.cloneNode(true);
    node.dataset.id = id;
    node.querySelector('.id').textContent = id;
    imagesList.prepend(node);
    startPolling(node, id);
    wireActions(node, id);
  }

  function wireActions(node, id) {
    const viewBtn = node.querySelector('.view');
    const delBtn = node.querySelector('.delete');

    viewBtn.addEventListener('click', async () => {
      try {
        const res = await fetch(API.image(id));
        if (!res.ok) {
          const j = await res.json().catch(() => ({}));
          throw new Error(j.message || 'Не удалось получить изображение');
        }
        const blob = await res.blob();
        const url = URL.createObjectURL(blob);
        const img = node.querySelector('.preview');
        img.src = url;
        img.classList.remove('hidden');
        node.querySelector('.placeholder').classList.add('hidden');
        showToast('Изображение получено');

        // Also trigger a download to the user's default Downloads folder
        const mime = blob.type || '';
        const ext = mime.includes('png') ? 'png' : mime.includes('gif') ? 'gif' : 'jpg';
        const a = document.createElement('a');
        a.href = url;
        a.download = `image_${id}.${ext}`;
        document.body.appendChild(a);
        a.click();
        a.remove();
      } catch (e) {
        showToast(e.message, 3500);
      }
    });

    delBtn.addEventListener('click', async () => {
      if (!confirm('Удалить изображение?')) return;
      try {
        const res = await fetch(API.delete(id), { method: 'DELETE' });
        const j = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(j.message || 'Не удалось удалить');
        node.remove();
        showToast('Изображение удалено');
      } catch (e) {
        showToast(e.message, 3500);
      }
    });
  }

  function startPolling(node, id) {
    const statusEl = node.querySelector('.status');
    const viewBtn = node.querySelector('.view');
    const delBtn = node.querySelector('.delete');
    const img = node.querySelector('.preview');
    const placeholder = node.querySelector('.placeholder');

    let pct = 10;
    let timer = setInterval(async () => {
      try {
        const res = await fetch(API.info(id));
        const json = await res.json();
        if (!res.ok) throw new Error(json.message || 'status fail');

        // API returns either string message or Image object
        let status = 'in_progress';
        if (typeof json.result === 'string') {
          // Could be "in processing, wait please" or error string
          status = json.result.includes('finished') ? 'finished' : 'in_progress';
        } else if (json.result && json.result.status) {
          status = json.result.status;
        }

        const map = statusMap[status] || ['in progress', false];
        statusEl.textContent = map[0];

        // fake progress bar advance while waiting
        const bar = node.querySelector('.progress .bar');
        pct = Math.min(95, pct + Math.random() * 10);
        if (bar) bar.style.width = pct + '%';

        if (map[1]) {
          clearInterval(timer);
          viewBtn.disabled = false;
          delBtn.disabled = false;
          if (bar) bar.style.width = '100%';

          // auto-fetch preview once
          try {
            const resImg = await fetch(API.image(id));
            if (resImg.ok) {
              const blob = await resImg.blob();
              const url = URL.createObjectURL(blob);
              img.src = url;
              img.classList.remove('hidden');
              placeholder.classList.add('hidden');
              showToast('Готово');
            }
          } catch {}
        } else {
          viewBtn.disabled = true;
          delBtn.disabled = true;
        }
      } catch (e) {
        // ignore transient errors
      }
    }, 1500);
  }
})();


