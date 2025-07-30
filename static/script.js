// ===== SOLUCIÓN INMEDIATA - DEFINIR FUNCIÓN GLOBAL AL INICIO =====
console.log('🚀 Definiendo convertFile inmediatamente...');

window.convertFile = function() {
    console.log('🔄 convertFile ejecutándose...');
    
    const fileInput = document.getElementById('txtFile');
    const statusElement = document.getElementById('converterStatusIndicator');
    
    // Validaciones básicas
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
        if (statusElement) {
            statusElement.className = 'alert alert-warning mb-3';
            statusElement.innerHTML = '<i class="fas fa-exclamation-triangle me-2"></i>Por favor selecciona un archivo primero';
        }
        return;
    }
    
    const file = fileInput.files[0];
    const fileName = file.name.toLowerCase();
    
    // Verificar extensión
    const supportedExtensions = ['.txt', '.json', '.md', '.csv', '.xml', '.yaml', '.yml'];
    const isSupported = supportedExtensions.some(ext => fileName.endsWith(ext));
    
    if (!isSupported) {
        if (statusElement) {
            statusElement.className = 'alert alert-danger mb-3';
            statusElement.innerHTML = '<i class="fas fa-times me-2"></i>Archivo no soportado. Use: .txt, .json, .md, .csv, .xml, .yaml, .yml';
        }
        return;
    }
    
    // Verificar tamaño
    if (file.size > 10 * 1024 * 1024) {
        if (statusElement) {
            statusElement.className = 'alert alert-danger mb-3';
            statusElement.innerHTML = '<i class="fas fa-times me-2"></i>Archivo muy grande (máximo 10MB)';
        }
        return;
    }
    
    // Mostrar estado de conversión
    if (statusElement) {
        statusElement.className = 'alert alert-info mb-3';
        statusElement.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>Convirtiendo archivo...';
    }
    
    // Crear FormData y enviar
    const formData = new FormData();
    formData.append('txtFile', file);
    
    fetch('/api/convert-to-go', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        return response.json();
    })
    .then(result => {
        console.log('✅ Conversión exitosa:', result);
        
        if (result.success) {
            // Mostrar éxito
            if (statusElement) {
                statusElement.className = 'alert alert-success mb-3';
                statusElement.innerHTML = '<i class="fas fa-check me-2"></i>¡Conversión exitosa!';
            }
            
            // Mostrar código generado
            showGeneratedCode(result);
        } else {
            // Mostrar error del servidor
            if (statusElement) {
                statusElement.className = 'alert alert-danger mb-3';
                statusElement.innerHTML = '<i class="fas fa-times me-2"></i>' + (result.error || 'Error en la conversión');
            }
        }
    })
    .catch(error => {
        console.error('❌ Error:', error);
        if (statusElement) {
            statusElement.className = 'alert alert-danger mb-3';
            statusElement.innerHTML = '<i class="fas fa-times me-2"></i>Error: ' + error.message;
        }
    });
};

// ===== FUNCIÓN PARA MOSTRAR CÓDIGO GENERADO =====
function showGeneratedCode(result) {
    const contentElement = document.getElementById('generatedCodeContent');
    if (!contentElement) return;
    
    const goCode = result.go_code || 'Error: No se generó código';
    const filename = result.download_filename || 'generated_code.go';
    
    contentElement.innerHTML = `
        <div class="alert alert-success mb-3">
            <h6><i class="fas fa-check-circle me-2"></i>¡Archivo convertido exitosamente!</h6>
            <p class="mb-0">Tu archivo ha sido convertido a código Go automáticamente.</p>
        </div>
        
        <div class="card border-primary mb-3">
            <div class="card-header bg-primary text-white">
                <div class="d-flex justify-content-between align-items-center">
                    <h6 class="mb-0"><i class="fab fa-golang me-2"></i>Código Go Generado</h6>
                    <div>
                        <button class="btn btn-light btn-sm me-2" onclick="copyCode()">
                            <i class="fas fa-copy me-1"></i>Copiar
                        </button>
                        <button class="btn btn-success btn-sm" onclick="downloadCode('${filename}')">
                            <i class="fas fa-download me-1"></i>Descargar
                        </button>
                    </div>
                </div>
            </div>
            <div class="card-body p-0">
                <pre id="generatedCode" style="margin: 0; padding: 15px; background: #f8f9fa; max-height: 400px; overflow-y: auto; font-size: 14px;">${escapeHtml(goCode)}</pre>
            </div>
        </div>
        
        <div class="row g-3">
            <div class="col-md-6">
                <div class="card border-info">
                    <div class="card-header bg-info text-white">
                        <h6 class="mb-0">Información</h6>
                    </div>
                    <div class="card-body">
                        <small><strong>Archivo:</strong> ${result.original_file || 'archivo'}</small><br>
                        <small><strong>Tamaño:</strong> ${formatBytes(result.file_size || 0)}</small><br>
                        <small><strong>Package:</strong> main (automático)</small><br>
                        <small><strong>Variable:</strong> textContent (automático)</small>
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="card border-success">
                    <div class="card-header bg-success text-white">
                        <h6 class="mb-0">Siguiente Paso</h6>
                    </div>
                    <div class="card-body">
                        <button class="btn btn-success btn-sm w-100 mb-2" onclick="copyCode()">
                            <i class="fas fa-copy me-2"></i>Copiar Código
                        </button>
                        <button class="btn btn-primary btn-sm w-100 mb-2" onclick="downloadCode('${filename}')">
                            <i class="fas fa-download me-2"></i>Descargar ${filename}
                        </button>
                        <button class="btn btn-outline-secondary btn-sm w-100" onclick="resetForm()">
                            <i class="fas fa-refresh me-2"></i>Nuevo Archivo
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// ===== FUNCIONES AUXILIARES =====
window.copyCode = function() {
    const codeElement = document.getElementById('generatedCode');
    if (codeElement) {
        navigator.clipboard.writeText(codeElement.textContent)
            .then(() => {
                alert('Código copiado al portapapeles');
            })
            .catch(() => {
                alert('No se pudo copiar automáticamente. Selecciona el texto manualmente.');
            });
    }
};

window.downloadCode = function(filename) {
    const codeElement = document.getElementById('generatedCode');
    if (codeElement) {
        const blob = new Blob([codeElement.textContent], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename || 'generated_code.go';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        alert('Archivo descargado');
    }
};

window.resetForm = function() {
    const form = document.getElementById('txtUploadForm');
    const content = document.getElementById('generatedCodeContent');
    const status = document.getElementById('converterStatusIndicator');
    
    if (form) form.reset();
    if (status) {
        status.className = 'alert alert-info mb-3';
        status.innerHTML = '<i class="fas fa-upload me-2"></i>Sube un archivo para convertir...';
    }
    if (content) {
        content.innerHTML = `
            <div class="text-center py-4">
                <i class="fab fa-golang fa-3x text-muted mb-3"></i>
                <h5 class="text-muted mb-3">Conversor Simplificado</h5>
                <p class="text-muted">Sube cualquier archivo de texto y se convertirá automáticamente a código Go.</p>
            </div>
        `;
    }
};

// ===== FUNCIONES DEL PARSER JSON =====
window.parseJSON = function() {
    const input = document.getElementById('jsonInput');
    if (!input) return;
    
    const jsonText = input.value.trim();
    if (!jsonText) {
        alert('Por favor ingresa un JSON');
        return;
    }
    
    fetch('/api/parse', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ json: jsonText })
    })
    .then(response => response.json())
    .then(result => {
        const resultElement = document.getElementById('resultContent');
        if (resultElement) {
            if (result.success) {
                resultElement.innerHTML = `
                    <div class="alert alert-success">
                        <h6>✅ JSON parseado exitosamente</h6>
                    </div>
                    <div class="card">
                        <div class="card-header">Resultado</div>
                        <div class="card-body">
                            <pre style="background: #f8f9fa; padding: 10px; border-radius: 5px;">${JSON.stringify(result.result, null, 2)}</pre>
                        </div>
                    </div>
                `;
            } else {
                resultElement.innerHTML = `
                    <div class="alert alert-danger">
                        <h6>❌ Error de parseo</h6>
                        <p class="mb-0">${result.error}</p>
                    </div>
                `;
            }
        }
    })
    .catch(error => {
        console.error('Error:', error);
        const resultElement = document.getElementById('resultContent');
        if (resultElement) {
            resultElement.innerHTML = `
                <div class="alert alert-danger">
                    <h6>❌ Error de conexión</h6>
                    <p class="mb-0">${error.message}</p>
                </div>
            `;
        }
    });
};

window.clearInput = function() {
    const input = document.getElementById('jsonInput');
    const result = document.getElementById('resultContent');
    
    if (input) input.value = '';
    if (result) {
        result.innerHTML = `
            <div class="text-center py-4">
                <i class="fas fa-trash-alt fa-3x text-muted mb-3"></i>
                <h5 class="text-muted">Área limpiada</h5>
                <p class="text-muted">Ingresa un nuevo JSON para parsear.</p>
            </div>
        `;
    }
    updateStats();
};

window.formatJSON = function() {
    const input = document.getElementById('jsonInput');
    if (!input) return;
    
    try {
        const parsed = JSON.parse(input.value);
        input.value = JSON.stringify(parsed, null, 2);
        updateStats();
        alert('JSON formateado correctamente');
    } catch (error) {
        alert('JSON inválido, no se puede formatear');
    }
};

window.toggleTheme = function() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
};

// ===== FUNCIONES AUXILIARES =====
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function updateStats() {
    const input = document.getElementById('jsonInput');
    const charCount = document.getElementById('charCount');
    const lineCount = document.getElementById('lineCount');
    
    if (input && charCount && lineCount) {
        const value = input.value;
        charCount.textContent = value.length.toLocaleString();
        lineCount.textContent = ((value.match(/\n/g) || []).length + 1).toLocaleString();
    }
}

function updateThemeIcon(theme) {
    const icon = document.getElementById('themeIcon');
    const text = document.getElementById('themeText');
    
    if (icon && text) {
        if (theme === 'dark') {
            icon.className = 'fas fa-sun';
            text.textContent = 'Modo Día';
        } else {
            icon.className = 'fas fa-moon';
            text.textContent = 'Modo Noche';
        }
    }
}

// ===== INICIALIZACIÓN =====
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 DOM cargado, inicializando...');
    
    // Cargar tema
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
    
    // Configurar input JSON
    const jsonInput = document.getElementById('jsonInput');
    if (jsonInput) {
        jsonInput.addEventListener('input', updateStats);
    }
    
    // Actualizar estadísticas
    updateStats();
    
    // Cargar ejemplos
    loadExamples();
    
    console.log('✅ Inicialización completa');
});

async function loadExamples() {
    try {
        const response = await fetch('/api/examples');
        if (response.ok) {
            const examples = await response.json();
            console.log('✅ Ejemplos cargados:', examples);
            renderExamples(examples);
        }
    } catch (error) {
        console.log('⚠️ Error cargando ejemplos:', error);
        renderDefaultExamples();
    }
}

function renderExamples(examples) {
    const grid = document.getElementById('examplesGrid');
    if (!grid || !examples.ejemplos) return;
    
    grid.innerHTML = '';
    examples.ejemplos.slice(0, 6).forEach((example, index) => {
        const col = document.createElement('div');
        col.className = 'col-md-6 col-lg-4 mb-2';
        
        const btn = document.createElement('button');
        btn.className = 'btn btn-outline-primary btn-sm w-100';
        btn.innerHTML = `<i class="fas fa-play me-1"></i>${example.nombre}`;
        btn.onclick = () => {
            const input = document.getElementById('jsonInput');
            if (input) {
                input.value = example.json;
                updateStats();
                setTimeout(() => window.parseJSON(), 500);
            }
        };
        
        col.appendChild(btn);
        grid.appendChild(col);
    });
}

function renderDefaultExamples() {
    const grid = document.getElementById('examplesGrid');
    if (!grid) return;
    
    const defaultExamples = [
        { nombre: 'Objeto Simple', json: '{"name": "Juan", "age": 30}' },
        { nombre: 'Array', json: '["a", "b", "c"]' },
        { nombre: 'Anidado', json: '{"user": {"name": "Ana"}}' }
    ];
    
    grid.innerHTML = '';
    defaultExamples.forEach(example => {
        const col = document.createElement('div');
        col.className = 'col-md-6 col-lg-4 mb-2';
        
        const btn = document.createElement('button');
        btn.className = 'btn btn-outline-primary btn-sm w-100';
        btn.innerHTML = `<i class="fas fa-play me-1"></i>${example.nombre}`;
        btn.onclick = () => {
            const input = document.getElementById('jsonInput');
            if (input) {
                input.value = example.json;
                updateStats();
                setTimeout(() => window.parseJSON(), 500);
            }
        };
        
        col.appendChild(btn);
        grid.appendChild(col);
    });
}

// ===== VERIFICACIÓN FINAL =====
console.log('✅ Script cargado completamente');
console.log('✅ convertFile disponible:', typeof window.convertFile === 'function');
console.log('✅ parseJSON disponible:', typeof window.parseJSON === 'function');
console.log('🎯 Listo para usar - el botón debería funcionar ahora');

// ===== ATAJOS DE TECLADO =====
document.addEventListener('keydown', function(e) {
    if (e.ctrlKey && e.key === 'Enter') {
        e.preventDefault();
        window.parseJSON();
    }
});