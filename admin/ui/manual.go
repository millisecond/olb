package ui

import (
	"html/template"
	"net/http"
)

type ManualHandler struct {
	Color    string
	Title    string
	Version  string
	Commands string
}

func (h *ManualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmplManual.ExecuteTemplate(w, "manual", h)
}

var tmplManual = template.Must(template.New("manual").Parse(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>fabio{{if .Title}} - {{.Title}}{{end}}</title>
	<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.97.3/css/materialize.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.97.3/js/materialize.min.js"></script>
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>

	<style type="text/css">
		.footer { padding-top: 10px; }
		.logo { height: 32px; margin: 0 auto; display: block; }
	</style>
</head>
<body>

<nav class="top-nav {{.Color}}">

	<div class="container">
		<div class="nav-wrapper">
			<a href="/" class="brand-logo">fabio{{if .Title}} - {{.Title}}{{end}}</a>
			<ul id="nav-mobile" class="right hide-on-med-and-down">
				<li><a href="/routes">Routes</a></li>
				<li><a href="https://github.com/fabiolb/fabio/blob/master/CHANGELOG.md">{{.Version}}</a></li>
				<li><a href="https://github.com/fabiolb/fabio">Github</a></li>
			</ul>
		</div>
	</div>

</nav>

<div class="container">

	<div class="section">
		<h5>Manual Overrides</h5>

		<div class="row">
			<form class="col s12">
				<input type="hidden" name="version">
				<div class="row">
					<div class="input-field col s12">
						<textarea id="textarea1" class="materialize-textarea"></textarea>
						<label for="textarea1"></label>
					</div>
				</div>
			</form>
			<button class="btn waves-effect waves-light" name="save">Save</button>
			<button class="btn waves-effect waves-light" name="help">Help</button>
		</div>

		<div class="row">
			<pre class="help hide">{{.Commands}}</pre>
		</div>
	</div>

	<div class="section footer">
		<img class="logo" src="/logo.svg">
	</div>

</div>

<script>
$(function(){
	var params={};window.location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,function(str,key,value){params[key] = value;});

	$.get("/api/manual", function(data) {
		$("input[name=version]").val(data.version);
		$("textarea>label").val("Version " + data.version);
		$("#textarea1").val(data.value);
		$("#textarea1").trigger('autoresize');
	});

	$("button[name=help]").click(function() {
		$("pre.help").toggleClass("hide");
	});

	$("button[name=save]").click(function() {
		var data = {
			value   : $("#textarea1").val(),
			version : $("input[name=version]").val()
		}
		$.ajax('/api/manual', {
			type: 'PUT',
			data: JSON.stringify(data),
			contentType: 'application/json',
			statusCode: {
				400: function(jqXHR, textStatus, err) { alert(err); },
				409: function(jqXHR, textStatus, err) { alert(err); },
				500: function(jqXHR, textStatus, err) { alert(err); }
			},
			success: function() {
				window.location.reload();
			}
		});
	});
})
</script>

</body>
</html>
`))
