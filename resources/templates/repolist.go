package templates

// TODO: Stable access order of Hooks map
// TODO: Map .ToLower func into this template to lowercase the $key in the URL
const RepoList = `
	{{ define "content" }}
	<br/><br/>
	<div>
	{{ range . }}
	{{ $repopath := .FullName }}
	<div><b><a href=/repos/{{ $repopath }}/hooks>{{ $repopath }}</a></b></div>
	<div><b>Active validators</b>:<br>
	{{ range $key, $value := .Hooks }}
		{{ if $value }}
		{{ $key }}: <a href="/results/{{ $key }}/{{ $repopath }}">results</a><br>
		{{ end }}
	{{ end }}
	</div>
	<div>{{.Description}} {{.Website}}</div>
	<hr>
	{{ end }}
	</div>
	{{ end }}
`

// TODO: Map .ToLower func into this template to lowercase the $key in the URL
const RepoPage = `
	{{ define "content" }}
	<br/><br/>
	<div>
	<div><b>{{.FullName}}</b></div>
	<div><b>Available validators</b>:<br>
	{{ range $key, $value := .Hooks }}
		{{ $key }}
		{{ if $value }}
		[Enabled] <a href="/repos/{{$.FullName}}/{{ $key }}/disable">disable</a>
		{{ else }}
		[Disabled] <a href="/repos/{{$.FullName}}/{{ $key }}/enable">enable</a>
		{{ end }}
		<br>
	{{ end }}
	</div>
	<div>{{.Description}} {{.Website}}</div>
	<hr>
	</div>
	{{ end }}
`
