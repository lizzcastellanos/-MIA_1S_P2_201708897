import React, { useState, useRef } from 'react';
import "../Stylesheets/navegador.css";

const CommandInterface = ({ apiEndpoint = "http://localhost:8080/analizar" }) => {
    const [inputCode, setInputCode] = useState('');
    const [outputResult, setOutputResult] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const fileInputRef = useRef(null);

    const handleCodeChange = (e) => setInputCode(e.target.value);
    const clearAll = () => {
        setInputCode('');
        setOutputResult('');
    };

    const executeCommand = async (e) => {
        e.preventDefault();
        if (!inputCode.trim()) return;

        setIsLoading(true);
        
        try {
            const response = await fetch(apiEndpoint, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ text: inputCode })
            });

            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            
            const { message } = await response.json();
            setOutputResult(message);
        } catch (error) {
            setOutputResult(`Error: ${error.message}`);
            console.error('API Request failed:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleFileUpload = () => fileInputRef.current?.click();

    const processFile = (e) => {
        const file = e.target.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (e) => setInputCode(e.target.result);
        reader.readAsText(file);
    };

    return (
        <div className="command-panel">
            <div className="io-container">
                <div className="code-section">
                    <label htmlFor="code-input" className="section-label">
                        <i className="icon fas fa-code"></i> EDITOR
                    </label>
                    <textarea
                        id="code-input"
                        className="code-editor"
                        value={inputCode}
                        onChange={handleCodeChange}
                        placeholder="// Ingresa tus comandos aquí..."
                        spellCheck="false"
                    />
                </div>

                <div className="result-section">
                    <label htmlFor="result-output" className="section-label">
                        <i className="icon fas fa-terminal"></i> RESULTADO
                    </label>
                    <textarea
                        id="result-output"
                        className="result-display"
                        value={outputResult}
                        readOnly
                    />
                </div>
            </div>

            <div className="action-buttons">
                <button 
                    onClick={executeCommand}
                    disabled={isLoading || !inputCode.trim()}
                    className={`execute-btn ${isLoading ? 'loading' : ''}`}
                >
                    {isLoading ? (
                        <>
                            <span className="spinner"></span> Procesando...
                        </>
                    ) : (
                        'Ejecutar'
                    )}
                </button>
                
                <button onClick={clearAll} className="secondary-btn">
                    <i className="fas fa-broom"></i> Limpiar
                </button>
                
                <div className="file-upload-wrapper">
                    <button onClick={handleFileUpload} className="secondary-btn">
                        <i className="fas fa-folder-open"></i> Abrir Archivo
                    </button>
                    <input
                        type="file"
                        ref={fileInputRef}
                        onChange={processFile}
                        accept=".txt,.js,.py,.java"
                        style={{ display: 'none' }}
                    />
                </div>
            </div>
        </div>
    );
};

export default CommandInterface;