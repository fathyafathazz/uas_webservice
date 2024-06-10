import SwaggerUIBundle from 'https://cdn.skypack.dev/swagger-ui-dist/swagger-ui-bundle.js';
import SwaggerUIStandalonePreset from 'https://cdn.skypack.dev/swagger-ui-dist/swagger-ui-standalone-preset.js';

export let URLData = "./openapi.yaml";

export const UIData = SwaggerUIBundle({
    url: URLData, //Location of Open API spec in the repo
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
    ],
    plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
    ]
});

export function setSwagger() {
    return UIData
}