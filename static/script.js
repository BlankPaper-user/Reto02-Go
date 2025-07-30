let examples = [];

// ===== FUNCIONES GLOBALES (DISPONIBLES DESDE HTML) =====
window.convertFile = function() {
    console.log('🔄 Función convertFile() llamada globalmente');
    const form = document.getElementById('txtUploadForm');
    if (form) {
        const event = new Event('submit', { bubbles: true, cancelable: true });
        handleTxtConversion(event);
    }
};

// Hacer todas las funciones principales disponibles globalmente
window.parseJSON = parseJSON;
window.clearInput = clearInput;  
window.formatJSON = formatJSON;
window.showExamplesModal = showExamplesModal;
window.selectExampleFromModal = selectExampleFromModal;
window.toggleTheme = toggleTheme;
window.copyGoCode = copyGoCode;
window.downloadGoCode = downloadGoCode;
window.resetConverter = resetConverter;

// ===== INICIALIZACIÓN =====
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 Inicializando aplicación...');
    
    loadExamples();
    updateStats();
    loadTheme();
    
    // Actualizar estadísticas cuando el usuario escriba
    document.getElementById('jsonInput').addEventListener('input', updateStats);
    
    // IMPORTANTE: Configurar el conversor TXT → Go
    setupTxtConverter();
    
    console.log('✅ Aplicación inicializada correctamente');
    console.log('✅ Función convertFile disponible globalmente:', typeof window.convertFile);
});

// ===== CONFIGURACIÓN DEL CONVERSOR TXT → Go =====
function setupTxtConverter() {
    const txtUploadForm = document.getElementById('txtUploadForm');
    if (txtUploadForm) {
        console.log('✅ Configurando conversor TXT → Go...');
        
        // Prevenir el envío normal del formulario
        txtUploadForm.addEventListener('submit', function(event) {
            event.preventDefault();
            event.stopPropagation();
            console.log('🔄 Formulario interceptado, iniciando conversión...');
            handleTxtConversion(event);
        });
        
        console.log('✅ Conversor TXT → Go configurado correctamente');
    } else {
        console.error('❌ No se encontró el formulario txtUploadForm');
    }
}

// ===== FUNCIONES DEL TEMA =====
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
}

function loadTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
}

function updateThemeIcon(theme) {
    const icon = document.getElementById('themeIcon');
    const text = document.getElementById('themeText');
    
    if (theme === 'dark') {
        icon.className = 'fas fa-sun';
        text.textContent = 'Modo Día';
    } else {
        icon.className = 'fas fa-moon';
        text.textContent = 'Modo Noche';
    }
}

// ===== FUNCIONES DE ESTADÍSTICAS =====
function updateStats() {
    const input = document.getElementById('jsonInput').value;
    document.getElementById('charCount').textContent = input.length.toLocaleString();
    document.getElementById('lineCount').textContent = ((input.match(/\n/g) || []).length + 1).toLocaleString();
}

// ===== GESTIÓN DE EJEMPLOS JSON =====
async function loadExamples() {
    try {
        console.log('🔄 Cargando ejemplos desde el servidor...');
        updateStatus('info', '<i class="fas fa-spinner fa-spin me-2"></i>Cargando ejemplos...');
        
        const response = await fetch('/api/examples');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('✅ Ejemplos cargados desde servidor:', data);
        examples = data;
        renderExamples();
        
        updateStatus('success', '<i class="fas fa-check me-2"></i>Ejemplos cargados correctamente');
        setTimeout(() => {
            updateStatus('info', '<i class="fas fa-clock me-2"></i>Esperando entrada...');
        }, 2000);
        
    } catch (error) {
        console.error('❌ Error cargando ejemplos:', error);
        updateStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>Usando ejemplos por defecto');
        renderDefaultExamples();
        
        setTimeout(() => {
            updateStatus('info', '<i class="fas fa-clock me-2"></i>Esperando entrada...');
        }, 2000);
    }
}

function showExamplesModal() {
    console.log('📚 Mostrando modal de ejemplos');
    populateExamplesModal();
}

function populateExamplesModal() {
    const validContainer = document.getElementById('validExamples');
    const invalidContainer = document.getElementById('invalidExamples');
    
    validContainer.innerHTML = '';
    invalidContainer.innerHTML = '';

    // Mostrar ejemplos válidos
    if (examples.ejemplos && examples.ejemplos.length > 0) {
        examples.ejemplos.forEach((example, index) => {
            const card = createExampleCard(example, index, 'valid');
            validContainer.appendChild(card);
        });
    } else {
        validContainer.innerHTML = '<div class="alert alert-info"><i class="fas fa-info-circle me-2"></i>No hay ejemplos válidos disponibles.</div>';
    }

    // Mostrar ejemplos inválidos
    if (examples.ejemplos_invalidos && examples.ejemplos_invalidos.length > 0) {
        examples.ejemplos_invalidos.forEach((example, index) => {
            const card = createExampleCard(example, index, 'invalid');
            invalidContainer.appendChild(card);
        });
    } else {
        invalidContainer.innerHTML = '<div class="alert alert-warning"><i class="fas fa-exclamation-triangle me-2"></i>No hay ejemplos inválidos disponibles.</div>';
    }
}

function createExampleCard(example, index, type) {
    const card = document.createElement('div');
    card.className = 'card mb-3';
    
    const borderClass = type === 'valid' ? 'border-success' : 'border-danger';
    card.classList.add(borderClass);
    
    const headerClass = type === 'valid' ? 'bg-success' : 'bg-danger';
    const icon = type === 'valid' ? 'fa-check' : 'fa-times';
    
    card.innerHTML = `
        <div class="card-header ${headerClass} text-white py-2">
            <div class="d-flex justify-content-between align-items-center">
                <small class="fw-bold"><i class="fas ${icon} me-1"></i>${example.nombre}</small>
                <button class="btn btn-sm btn-light" onclick="selectExampleFromModal('${type}', ${index})" data-bs-dismiss="modal">
                    <i class="fas fa-arrow-right me-1"></i>Usar
                </button>
            </div>
        </div>
        <div class="card-body py-2">
            <div class="code-container">
                <pre class="mb-0 p-3" style="font-size: 0.8em; max-height: 120px; overflow-y: auto;">${example.json}</pre>
            </div>
        </div>
    `;
    
    return card;
}

function selectExampleFromModal(type, index) {
    let selectedExample;
    
    if (type === 'valid') {
        selectedExample = examples.ejemplos[index];
    } else {
        selectedExample = examples.ejemplos_invalidos[index];
    }
    
    console.log(`📝 Seleccionando ejemplo ${type}:`, selectedExample);
    setExample(selectedExample.json);
    
    // Mostrar mensaje informativo
    if (type === 'invalid') {
        updateStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>¡Cuidado! Este es un ejemplo inválido para probar errores');
    } else {
        updateStatus('success', '<i class="fas fa-check me-2"></i>Ejemplo cargado - presiona "Parsear JSON"');
    }
}

function renderExamples() {
    const grid = document.getElementById('examplesGrid');
    grid.innerHTML = '';

    let examplesList = [];
    
    // Verificar diferentes estructuras de respuesta
    if (examples.ejemplos && Array.isArray(examples.ejemplos)) {
        examplesList = examples.ejemplos;
    } else if (Array.isArray(examples)) {
        examplesList = examples;
    } else {
        console.log('Estructura de ejemplos no reconocida:', examples);
        renderDefaultExamples();
        return;
    }

    console.log('Renderizando ejemplos:', examplesList);

    if (examplesList.length === 0) {
        renderDefaultExamples();
        return;
    }

    examplesList.slice(0, 6).forEach((example, index) => {
        const col = document.createElement('div');
        col.className = 'col-md-6 col-lg-4 mb-2';
        
        const btn = document.createElement('button');
        btn.className = 'btn example-btn btn-sm w-100';
        btn.innerHTML = `<i class="fas fa-play me-1"></i>${example.nombre || example.name || `Ejemplo ${index + 1}`}`;
        btn.onclick = () => {
            console.log('Cargando ejemplo:', example);
            setExample(example.json);
        };
        
        col.appendChild(btn);
        grid.appendChild(col);
    });
}

function renderDefaultExamples() {
    console.log('📋 Renderizando ejemplos por defecto');
    const defaultExamples = [
        { nombre: 'Objeto Simple', json: '{"name": "Juan", "age": 30}' },
        { nombre: 'Array', json: '["a", "b", "c"]' },
        { nombre: 'Anidado', json: '{"user": {"name": "Ana"}}' },
        { nombre: 'Mixto', json: '{"str": "texto", "num": 42, "bool": true}' }
    ];

    const grid = document.getElementById('examplesGrid');
    grid.innerHTML = '';

    defaultExamples.forEach((example, index) => {
        const col = document.createElement('div');
        col.className = 'col-md-6 col-lg-4 mb-2';
        
        const btn = document.createElement('button');
        btn.className = 'btn example-btn btn-sm w-100';
        btn.innerHTML = `<i class="fas fa-play me-1"></i>${example.nombre}`;
        btn.onclick = () => {
            console.log('🎯 Cargando ejemplo por defecto:', example);
            setExample(example.json);
        };
        
        col.appendChild(btn);
        grid.appendChild(col);
    });

    // Mensaje informativo
    const infoCol = document.createElement('div');
    infoCol.className = 'col-12 mt-2';
    infoCol.innerHTML = '<div class="alert alert-warning"><small><i class="fas fa-info-circle me-2"></i>Ejemplos por defecto (servidor no disponible)</small></div>';
    grid.appendChild(infoCol);
}

function setExample(jsonString) {
    console.log('Estableciendo ejemplo:', jsonString);
    document.getElementById('jsonInput').value = jsonString;
    updateStats();
    // Auto-parsear después de un breve delay
    setTimeout(() => parseJSON(), 500);
}

// ===== FUNCIONES DEL PARSER JSON =====
function clearInput() {
    document.getElementById('jsonInput').value = '';
    document.getElementById('resultContent').innerHTML = `
        <div class="text-center py-4">
            <i class="fas fa-trash-alt fa-3x text-muted mb-3"></i>
            <h5 class="text-muted mb-3">Área limpiada</h5>
            <p class="text-muted">Ingresa un nuevo JSON para parsear.</p>
        </div>
    `;
    document.getElementById('resultContainer').className = 'result-container';
    updateStatus('info', '<i class="fas fa-clock me-2"></i>Esperando entrada...');
    updateStats();
}

function formatJSON() {
    const input = document.getElementById('jsonInput').value.trim();
    if (!input) {
        updateStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>No hay contenido para formatear');
        return;
    }

    try {
        const parsed = JSON.parse(input);
        const formatted = JSON.stringify(parsed, null, 2);
        document.getElementById('jsonInput').value = formatted;
        updateStats();
        updateStatus('success', '<i class="fas fa-check me-2"></i>JSON formateado correctamente');
    } catch (error) {
        updateStatus('danger', '<i class="fas fa-times me-2"></i>No se puede formatear: JSON inválido');
    }
}

async function parseJSON() {
    const input = document.getElementById('jsonInput').value.trim();
    
    if (!input) {
        updateStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>Por favor ingresa un JSON');
        return;
    }

    // Mostrar estado de carga
    updateStatus('info', '<i class="fas fa-spinner fa-spin me-2"></i>Parseando...');

    try {
        console.log('Enviando JSON para parsear:', input);
        
        const response = await fetch('/api/parse', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ json: input })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        console.log('Resultado del parseo:', result);
        
        if (result.success) {
            displaySuccess(result.result);
        } else {
            displayError(result.error);
        }
    } catch (error) {
        console.error('Error en parseJSON:', error);
        displayError('Error de conexión: ' + error.message);
    }
}

function displaySuccess(result) {
    updateStatus('success', '<i class="fas fa-check me-2"></i>JSON parseado exitosamente');
    
    const container = document.getElementById('resultContainer');
    const content = document.getElementById('resultContent');
    
    container.className = 'result-container';
    
    // Formatear el resultado de manera legible
    const formattedResult = formatGoValue(result, 0);
    content.innerHTML = `
        <div class="alert alert-success mb-4" role="alert">
            <h6 class="alert-heading mb-2">
                <i class="fas fa-check-circle me-2"></i>¡Parseo exitoso!
            </h6>
            <p class="mb-0">El JSON ha sido procesado correctamente por el parser personalizado.</p>
        </div>

        <div class="row g-3">
            <div class="col-lg-6">
                <div class="card border-primary">
                    <div class="card-header bg-primary text-white">
                        <h6 class="mb-0">
                            <i class="fas fa-code me-2"></i>Resultado (estructura Go)
                        </h6>
                    </div>
                    <div class="card-body p-0">
                        <div class="code-container">
                            <pre class="mb-0 p-3" style="font-size: 0.85rem; max-height: 400px; overflow-y: auto;">${formattedResult}</pre>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="col-lg-6">
                <div class="card border-info">
                    <div class="card-header bg-info text-white">
                        <h6 class="mb-0">
                            <i class="fas fa-file-code me-2"></i>JSON formateado
                        </h6>
                    </div>
                    <div class="card-body p-0">
                        <div class="code-container">
                            <pre class="mb-0 p-3" style="font-size: 0.85rem; max-height: 400px; overflow-y: auto;">${JSON.stringify(result, null, 2)}</pre>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="alert alert-light border mt-4">
            <div class="row text-center g-3">
                <div class="col-md-4">
                    <i class="fas fa-check-circle text-success me-2"></i>
                    <small><strong>Sintaxis válida</strong></small>
                </div>
                <div class="col-md-4">
                    <i class="fas fa-cogs text-primary me-2"></i>
                    <small><strong>Parser personalizado</strong></small>
                </div>
                <div class="col-md-4">
                    <i class="fas fa-bolt text-warning me-2"></i>
                    <small><strong>Procesado en tiempo real</strong></small>
                </div>
            </div>
        </div>
    `;
}

function displayError(error) {
    updateStatus('danger', '<i class="fas fa-times me-2"></i>Error en el parseo');
    
    const container = document.getElementById('resultContainer');
    const content = document.getElementById('resultContent');
    
    container.className = 'result-container';
    content.innerHTML = `
        <div class="alert alert-danger mb-4" role="alert">
            <h6 class="alert-heading mb-2">
                <i class="fas fa-exclamation-triangle me-2"></i>Error de parseo detectado
            </h6>
            <p class="mb-0">El parser encontró un problema en la sintaxis del JSON.</p>
        </div>

        <div class="card border-danger mb-4">
            <div class="card-header bg-danger text-white">
                <h6 class="mb-0">
                    <i class="fas fa-bug me-2"></i>Detalle del error
                </h6>
            </div>
            <div class="card-body">
                <div class="code-container">
                    <pre class="mb-0 p-3 text-danger" style="font-size: 0.9rem; white-space: pre-wrap; background: #fff5f5;">${error}</pre>
                </div>
            </div>
        </div>

        <div class="card border-warning">
            <div class="card-header bg-warning">
                <h6 class="mb-0">
                    <i class="fas fa-lightbulb me-2"></i>Consejos para corregir
                </h6>
            </div>
            <div class="card-body">
                <div class="row g-3">
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>Verifica que todas las llaves <code>{}</code> y corchetes <code>[]</code> estén balanceados</small>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>Asegúrate de que las cadenas estén entre comillas dobles <code>"</code></small>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>No uses comas adicionales al final de objetos o arrays</small>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>Verifica la sintaxis de números y valores booleanos</small>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>Revisa los caracteres de escape en strings</small>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="d-flex align-items-start">
                            <i class="fas fa-check text-success me-2 mt-1"></i>
                            <small>Usa el formateo automático para detectar problemas</small>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="text-center mt-4">
            <button class="btn btn-outline-primary me-2" onclick="showExamplesModal()" data-bs-toggle="modal" data-bs-target="#examplesModal">
                <i class="fas fa-book me-2"></i>Ver ejemplos válidos
            </button>
            <button class="btn btn-outline-secondary" onclick="formatJSON()">
                <i class="fas fa-magic me-2"></i>Intentar formatear
            </button>
        </div>
    `;
}

function updateStatus(type, message) {
    const indicator = document.getElementById('statusIndicator');
    const alertClasses = {
        'info': 'alert-info',
        'success': 'alert-success', 
        'warning': 'alert-warning',
        'danger': 'alert-danger'
    };
    
    indicator.className = `alert ${alertClasses[type]} mb-3`;
    indicator.innerHTML = message;
}

function formatGoValue(value, indent = 0) {
    const spaces = '  '.repeat(indent);
    const nextSpaces = '  '.repeat(indent + 1);
    
    if (value === null) {
        return 'nil';
    } else if (typeof value === 'boolean') {
        return value.toString();
    } else if (typeof value === 'number') {
        return value.toString();
    } else if (typeof value === 'string') {
        return `"${value}"`;
    } else if (Array.isArray(value)) {
        if (value.length === 0) {
            return '[]interface{}{}';
        }
        let result = '[]interface{}{\n';
        value.forEach((item, index) => {
            result += nextSpaces + formatGoValue(item, indent + 1);
            if (index < value.length - 1) result += ',';
            result += '\n';
        });
        result += spaces + '}';
        return result;
    } else if (typeof value === 'object') {
        const keys = Object.keys(value);
        if (keys.length === 0) {
            return 'map[string]interface{}{}';
        }
        let result = 'map[string]interface{}{\n';
        keys.forEach((key, index) => {
            result += nextSpaces + `"${key}": ${formatGoValue(value[key], indent + 1)}`;
            if (index < keys.length - 1) result += ',';
            result += '\n';
        });
        result += spaces + '}';
        return result;
    }
    return String(value);
}

// ===== FUNCIONES DEL CONVERSOR ARCHIVO → Go =====

async function handleTxtConversion(event) {
    console.log('🚀 Iniciando handleTxtConversion...');
    
    // CRÍTICO: Prevenir cualquier envío del formulario
    if (event) {
        event.preventDefault();
        event.stopPropagation();
    }
    
    // Actualizar estado del conversor
    updateConverterStatus('info', '<i class="fas fa-spinner fa-spin me-2"></i>Convirtiendo archivo...');
    
    // Obtener datos del formulario
    const form = document.getElementById('txtUploadForm');
    const formData = new FormData(form);
    const fileInput = document.getElementById('txtFile');
    
    console.log('📋 Datos del formulario:', {
        packageName: formData.get('packageName'),
        variableName: formData.get('variableName'),
        conversionType: formData.get('conversionType'),
        file: fileInput.files[0]?.name
    });
    
    // Validaciones
    if (!fileInput.files || fileInput.files.length === 0) {
        updateConverterStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>Por favor selecciona un archivo');
        return false;
    }
    
    const file = fileInput.files[0];
    const fileName = file.name.toLowerCase();
    
    // SOPORTE PARA MÚLTIPLES TIPOS DE ARCHIVO
    const supportedExtensions = ['.txt', '.json', '.md', '.csv', '.xml', '.yaml', '.yml'];
    const isSupported = supportedExtensions.some(ext => fileName.endsWith(ext));
    
    if (!isSupported) {
        updateConverterStatus('danger', 
            '<i class="fas fa-times me-2"></i>Tipo de archivo no soportado. ' +
            'Archivos permitidos: .txt, .json, .md, .csv, .xml, .yaml, .yml'
        );
        return false;
    }
    
    if (file.size > 10 * 1024 * 1024) {
        updateConverterStatus('danger', '<i class="fas fa-times me-2"></i>El archivo es demasiado grande (máximo 10MB)');
        return false;
    }
    
    try {
        console.log('📤 Enviando petición al servidor...');
        
        // PETICIÓN AJAX CON FETCH
        const response = await fetch('/api/convert-to-go', {
            method: 'POST',
            body: formData,
        });
        
        console.log('📥 Respuesta del servidor:', response.status, response.statusText);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('✅ Conversión exitosa:', result);
        
        if (result.success) {
            displayGoCode(result);
            showFileInfo(result);
        } else {
            displayConversionError(result.error || 'Error desconocido en la conversión');
        }
        
    } catch (error) {
        console.error('❌ Error en conversión:', error);
        displayConversionError('Error de conexión: ' + error.message);
    }
    
    return false;
}

function updateConverterStatus(type, message) {
    const indicator = document.getElementById('converterStatusIndicator');
    if (!indicator) return;
    
    const alertClasses = {
        'info': 'alert-info',
        'success': 'alert-success', 
        'warning': 'alert-warning',
        'danger': 'alert-danger'
    };
    
    indicator.className = `alert ${alertClasses[type]} mb-3`;
    indicator.innerHTML = message;
}

function displayGoCode(result) {
    updateConverterStatus('success', '<i class="fas fa-check me-2"></i>Conversión exitosa');
    
    const container = document.getElementById('generatedCodeContainer');
    const content = document.getElementById('generatedCodeContent');
    
    if (!container || !content) return;
    
    container.className = 'result-container';
    
    // Usar concatenación de strings para evitar problemas con template literals
    content.innerHTML = 
        '<div class="alert alert-success mb-4" role="alert">' +
            '<h6 class="alert-heading mb-2">' +
                '<i class="fas fa-check-circle me-2"></i>¡Conversión exitosa!' +
            '</h6>' +
            '<p class="mb-0">El archivo ha sido convertido a código Go correctamente.</p>' +
        '</div>' +

        '<div class="card border-primary mb-4">' +
            '<div class="card-header bg-primary text-white">' +
                '<div class="d-flex justify-content-between align-items-center">' +
                    '<h6 class="mb-0">' +
                        '<i class="fab fa-golang me-2"></i>Código Go Generado' +
                    '</h6>' +
                    '<div>' +
                        '<button class="btn btn-light btn-sm me-2" onclick="copyGoCode()">' +
                            '<i class="fas fa-copy me-1"></i>Copiar' +
                        '</button>' +
                        '<button class="btn btn-success btn-sm" onclick="downloadGoCode(\'' + 
                        (result.download_filename || 'generated_code.go') + '\')">' +
                            '<i class="fas fa-download me-1"></i>Descargar' +
                        '</button>' +
                    '</div>' +
                '</div>' +
            '</div>' +
            '<div class="card-body p-0">' +
                '<div class="code-container">' +
                    '<pre id="generatedGoCode" class="mb-0 p-3" style="font-size: 0.85rem; max-height: 400px; overflow-y: auto;">' +
                    escapeHtml(result.go_code) +
                    '</pre>' +
                '</div>' +
            '</div>' +
        '</div>' +

        '<div class="row g-3">' +
            '<div class="col-md-6">' +
                '<div class="card border-info">' +
                    '<div class="card-header bg-info text-white">' +
                        '<h6 class="mb-0">' +
                            '<i class="fas fa-info-circle me-2"></i>Información de Conversión' +
                        '</h6>' +
                    '</div>' +
                    '<div class="card-body">' +
                        '<div class="row g-2">' +
                            '<div class="col-12">' +
                                '<small><strong>Archivo original:</strong> ' + (result.original_file || 'archivo') + '</small>' +
                            '</div>' +
                            '<div class="col-12">' +
                                '<small><strong>Tamaño:</strong> ' + formatBytes(result.file_size || 0) + '</small>' +
                            '</div>' +
                            '<div class="col-12">' +
                                '<small><strong>Tiempo:</strong> ' + (result.conversion_time || '0ms') + '</small>' +
                            '</div>' +
                            '<div class="col-12">' +
                                '<small><strong>Package:</strong> ' + (result.parameters?.package_name || 'main') + '</small>' +
                            '</div>' +
                            '<div class="col-12">' +
                                '<small><strong>Variable:</strong> ' + (result.parameters?.variable_name || 'textContent') + '</small>' +
                            '</div>' +
                            '<div class="col-12">' +
                                '<small><strong>Tipo:</strong> ' + (result.parameters?.conversion_type || 'variable') + '</small>' +
                            '</div>' +
                        '</div>' +
                    '</div>' +
                '</div>' +
            '</div>' +
            
            '<div class="col-md-6">' +
                '<div class="card border-success">' +
                    '<div class="card-header bg-success text-white">' +
                        '<h6 class="mb-0">' +
                            '<i class="fas fa-code me-2"></i>Siguiente Paso' +
                        '</h6>' +
                    '</div>' +
                    '<div class="card-body">' +
                        '<p class="mb-3"><small>Tu código Go está listo para usar:</small></p>' +
                        '<div class="d-grid gap-2">' +
                            '<button class="btn btn-success btn-sm" onclick="copyGoCode()">' +
                                '<i class="fas fa-copy me-2"></i>Copiar Código' +
                            '</button>' +
                            '<button class="btn btn-primary btn-sm" onclick="downloadGoCode(\'' + 
                            (result.download_filename || 'generated_code.go') + '\')">' +
                                '<i class="fas fa-download me-2"></i>Descargar ' + 
                                (result.download_filename || 'generated_code.go') +
                            '</button>' +
                            '<button class="btn btn-outline-secondary btn-sm" onclick="resetConverter()">' +
                                '<i class="fas fa-refresh me-2"></i>Nueva Conversión' +
                            '</button>' +
                        '</div>' +
                    '</div>' +
                '</div>' +
            '</div>' +
        '</div>';
}

function displayConversionError(error) {
    updateConverterStatus('danger', '<i class="fas fa-times me-2"></i>Error en la conversión');
    
    const container = document.getElementById('generatedCodeContainer');
    const content = document.getElementById('generatedCodeContent');
    
    if (!container || !content) return;
    
    container.className = 'result-container';
    content.innerHTML = 
        '<div class="alert alert-danger mb-4" role="alert">' +
            '<h6 class="alert-heading mb-2">' +
                '<i class="fas fa-exclamation-triangle me-2"></i>Error en la conversión' +
            '</h6>' +
            '<p class="mb-0">No se pudo convertir el archivo a código Go.</p>' +
        '</div>' +

        '<div class="card border-danger mb-4">' +
            '<div class="card-header bg-danger text-white">' +
                '<h6 class="mb-0">' +
                    '<i class="fas fa-bug me-2"></i>Detalle del error' +
                '</h6>' +
            '</div>' +
            '<div class="card-body">' +
                '<div class="code-container">' +
                    '<pre class="mb-0 p-3 text-danger" style="font-size: 0.9rem; white-space: pre-wrap; background: #fff5f5;">' +
                    error +
                    '</pre>' +
                '</div>' +
            '</div>' +
        '</div>' +

        '<div class="text-center">' +
            '<button class="btn btn-outline-primary" onclick="resetConverter()">' +
                '<i class="fas fa-refresh me-2"></i>Intentar de nuevo' +
            '</button>' +
        '</div>';
}

function showFileInfo(result) {
    const fileInfo = document.getElementById('fileInfo');
    const fileInfoContent = document.getElementById('fileInfoContent');
    
    if (!fileInfo || !fileInfoContent) return;
    
    fileInfoContent.innerHTML = 
        '<div class="col-md-6">' +
            '<small class="d-block"><strong>Nombre:</strong> ' + (result.original_file || 'archivo') + '</small>' +
        '</div>' +
        '<div class="col-md-6">' +
            '<small class="d-block"><strong>Tamaño:</strong> ' + formatBytes(result.file_size || 0) + '</small>' +
        '</div>' +
        '<div class="col-md-6">' +
            '<small class="d-block"><strong>Tiempo:</strong> ' + (result.conversion_time || '0ms') + '</small>' +
        '</div>' +
        '<div class="col-md-6">' +
            '<small class="d-block"><strong>Estado:</strong> <span class="text-success">Convertido</span></small>' +
        '</div>';
    
    fileInfo.style.display = 'block';
}

function copyGoCode() {
    const codeElement = document.getElementById('generatedGoCode');
    if (!codeElement) return;
    
    const text = codeElement.textContent;
    navigator.clipboard.writeText(text).then(() => {
        updateConverterStatus('success', '<i class="fas fa-check me-2"></i>Código copiado al portapapeles');
    }).catch(() => {
        updateConverterStatus('warning', '<i class="fas fa-exclamation-triangle me-2"></i>No se pudo copiar automáticamente');
    });
}

function downloadGoCode(filename) {
    const codeElement = document.getElementById('generatedGoCode');
    if (!codeElement) return;
    
    const text = codeElement.textContent;
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    
    const a = document.createElement('a');
    a.href = url;
    a.download = filename || 'generated_code.go';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    updateConverterStatus('success', '<i class="fas fa-check me-2"></i>Archivo descargado');
}

function resetConverter() {
    // Limpiar formulario
    document.getElementById('txtUploadForm').reset();
    
    // Ocultar información del archivo
    const fileInfo = document.getElementById('fileInfo');
    if (fileInfo) fileInfo.style.display = 'none';
    
    // Restaurar estado inicial
    const content = document.getElementById('generatedCodeContent');
    if (content) {
        content.innerHTML = 
            '<div class="text-center py-4">' +
                '<i class="fab fa-golang fa-3x text-muted mb-3"></i>' +
                '<h5 class="text-muted mb-3">Conversor Archivo → Go</h5>' +
                '<p class="text-muted mb-3">Sube un archivo y convierte su contenido a código Go.</p>' +
                
                '<div class="divider"></div>' +
                
                '<div class="row g-3">' +
                    '<div class="col-md-6">' +
                        '<div class="feature-item">' +
                            '<i class="fas fa-file-code text-primary mb-2"></i>' +
                            '<div><small><strong>TXT, JSON, MD</strong></small></div>' +
                        '</div>' +
                    '</div>' +
                    '<div class="col-md-6">' +
                        '<div class="feature-item">' +
                            '<i class="fas fa-function text-primary mb-2"></i>' +
                            '<div><small><strong>CSV, XML, YAML</strong></small></div>' +
                        '</div>' +
                    '</div>' +
                    '<div class="col-md-6">' +
                        '<div class="feature-item">' +
                            '<i class="fas fa-sitemap text-primary mb-2"></i>' +
                            '<div><small><strong>Estructuras</strong></small></div>' +
                        '</div>' +
                    '</div>' +
                    '<div class="col-md-6">' +
                        '<div class="feature-item">' +
                            '<i class="fas fa-list text-primary mb-2"></i>' +
                            '<div><small><strong>Variables y Funciones</strong></small></div>' +
                        '</div>' +
                    '</div>' +
                '</div>' +
            '</div>';
    }
    
    updateConverterStatus('info', '<i class="fas fa-upload me-2"></i>Sube un archivo para generar código Go...');
}

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

// ===== ATAJOS DE TECLADO =====

document.addEventListener('keydown', function(e) {
    if (e.ctrlKey && e.key === 'Enter') {
        e.preventDefault();
        parseJSON();
    } else if (e.ctrlKey && e.key === 'l') {
        e.preventDefault();
        clearInput();
    } else if (e.ctrlKey && e.shiftKey && e.key === 'F') {
        e.preventDefault();
        formatJSON();
    }
});

// ===== LOG DE INICIALIZACIÓN =====
console.log('✅ Script completo cargado correctamente');
console.log('🔧 Funciones globales disponibles:', Object.keys(window).filter(key => 
    ['convertFile', 'parseJSON', 'clearInput', 'formatJSON', 'showExamplesModal', 'toggleTheme'].includes(key)
));