package swagger

const indexTmpl string = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta
            name="description"
            content="SwaggerUI"
    />
    <title>{{.Title}}</title>
    <link href="https://fonts.cdnfonts.com/css/chicagoflf" rel="stylesheet">
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/capiorg/swagger-ui@0.0.1/static/css.min.css">
    {{- if .CustomStyle}}
    <style>
        {{.CustomStyle}}
    </style>
    {{- end}}
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.18.2/swagger-ui-bundle.js" crossorigin></script>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.18.2/swagger-ui-standalone-preset.js" crossorigin></script>
<script>
    window.onload = function() {
        config = {{.}};
        config.dom_id = '#swagger-ui';
        config.plugins = [
            {{- range $plugin := .Plugins }}
        {{$plugin}},
        {{- end}}
    ];
        config.presets = [
            {{- range $preset := .Presets }}
        {{$preset}},
        {{- end}}
    ];
        config.filter = {{.Filter.Value}}
        config.syntaxHighlight = {{.SyntaxHighlight.Value}}
        {{if .TagsSorter}}
        config.tagsSorter = {{.TagsSorter}}
        {{end}}
        {{if .OnComplete}}
        config.onComplete = {{.OnComplete}}
        {{end}}
        {{if .RequestInterceptor}}
        config.requestInterceptor = {{.RequestInterceptor}}
        {{end}}
        {{if .ResponseInterceptor}}
        config.responseInterceptor = {{.ResponseInterceptor}}
        {{end}}
        {{if .ModelPropertyMacro}}
        config.modelPropertyMacro = {{.ModelPropertyMacro}}
        {{end}}
        {{if .ParameterMacro}}
        config.parameterMacro = {{.ParameterMacro}}
        {{end}}
        const ui = SwaggerUIBundle(config);

        {{if .OAuth}}
        ui.initOAuth({{.OAuth}});
        {{end}}
        {{if .PreauthorizeBasic}}
        ui.preauthorizeBasic({{.PreauthorizeBasic}});
        {{end}}
        {{if .PreauthorizeApiKey}}
        ui.preauthorizeApiKey({{.PreauthorizeApiKey}});
        {{end}}

        window.ui = ui
    }
</script>
</body>
</html>
`
