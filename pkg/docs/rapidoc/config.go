package rapidoc

type (
	RenderStyle string
	SchemaStyle string
	ThemeStyle  string
)

const (
	RenderStyle_Read  RenderStyle = "read"
	RenderStyle_View  RenderStyle = "view"
	RenderStyle_Focus RenderStyle = "focused"

	SchemaStyle_Tree  SchemaStyle = "tree"
	SchemaStyle_Table SchemaStyle = "table"

	Theme_Dark  ThemeStyle = "dark"
	Theme_Light ThemeStyle = "light"
)

type RapiDocConfig struct {
	Title       string      `json:"tiltle,omitempty"`
	SpecURL     string      `json:"spec_url,omitempty"`
	HeaderText  string      `json:"header_text,omitempty"`
	LogoURL     string      `json:"logo_url,omitempty"`
	RenderStyle RenderStyle `json:"render_style,omitempty"`
	SchemaStyle SchemaStyle `json:"schema_style,omitempty"`
	Theme       ThemeStyle  `json:"theme,omitempty"`
}

func GetDefaultRapiDocConfig() RapiDocConfig {
	return RapiDocConfig{
		Title:       "API Documentation",
		SpecURL:     "./swagger.json",
		HeaderText:  "API Documentation",
		LogoURL:     "https://mrin9.github.io/RapiDoc/images/logo.png",
		RenderStyle: RenderStyle_Read,
		SchemaStyle: SchemaStyle_Tree,
		Theme:       Theme_Dark,
	}
}

func HtmlTemplateRapiDoc() string {
	return `<!doctype html>
	<html>	
	<head>
		<title>{{$.Title}}</title>
		<meta charset="utf-8">
		<link href="https://fonts.googleapis.com/css?family=Nunito" rel="stylesheet">
		<link href="https://fonts.googleapis.com/css2?family=Sarabun&display=swap" rel="stylesheet">
		<link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;600&family=Roboto+Mono&display=swap" rel="stylesheet">
		<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
		<style>
			rapi-doc::part(section-servers) { /* <<< targets the server div */
			background: #6b5b95;
			color:#d1c2e4;
			margin:0 24px 0 24px;
			border-radius: 5px;
		}
			rapi-doc::part(label-selected-server) { /* <<< targets selected server label */
			color: #fff;
		}
		</style>
	</head>
	
	<body>
		<rapi-doc 
		spec-url="{{$.SpecURL}}" 
		heading-text="{{$.HeaderText}}" 
		theme="{{$.Theme}}"
		show-header=false
		use-path-in-nav-bar=true
		show-method-in-nav-bar="as-colored-text"
		nav-bg-color = "#423368"
		primary-color ="#df75c4"
		bg-color = "#fae4f5"
		regular-font = "Nunito"
		mono-font="'Roboto Mono'" 
		render-style="{{$.RenderStyle}}"
		persist-auth="true"
		schema-style="{{$.SchemaStyle}}">
		<div slot="nav-logo" style="display: flex; align-items: center; justify-content: center;">
			<img src="{{$.LogoURL}}" style="width:75px">
		</div>
		</rapi-doc>
	</body>
	</html>`
}
