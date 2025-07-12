let examples = [];

// Cargar ejemplos al inicio
document.addEventListener('DOMContentLoaded', function() {
    loadExamples();
    updateStats();
    loadTheme();
    
    // Actualizar estadísticas cuando el usuario escriba
    document.getElementById('jsonInput').addEventListener('input', updateStats);
});

// Funciones del tema
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

function updateStats() {
    const input = document.getElementById('jsonInput').value;
    document.getElementById('charCount').textContent = input.length.toLocaleString();
    document.getElementById('lineCount').textContent = ((input.match(/\n/g) || []).length + 1).toLocaleString();
}

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

// Atajos de teclado
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