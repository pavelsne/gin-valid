package templates

// TODO: Stable access order of Hooks map
const RepoList = `
{{define "content"}}

<div class="explore repositories">
	<div class="ui container repository list">
		{{range .}}
			{{$repopath := .FullName}}
			<div class="item">
				<div class="ui grid">
					<div class="ui two wide column middle aligned center">
						<i class="mega-octicon octicon-repo"></i>
					</div>
					<div class="ui fourteen wide column">
						<div class="ui header">
							<a class="name" href="/repos/{{$repopath}}/hooks">{{$repopath}}</a>
							<div class="ui links text normal small">
								<span>Active validators</span>
						{{range $hookname, $hook := .Hooks}}
							{{if eq $hook.State 0}}
								<span> | {{$hookname | ToUpper}}: <a href="/results/{{$hookname | ToLower}}/{{$repopath}}">results</a> </span>
							{{end}}
						{{end}}
							</div>
						</div>
						<p class="has-emoji">{{.Description}}</p>
						<a href="{{.HTMLURL}}">Repository on GIN</a> | <a href="{{.HTMLURL}}/settings/hooks">Repository hooks</a>
					</div>
				</div>
			</div>
		{{end}}
	</div>
</div>
	{{end}}
`

// TODO: Map .ToLower func into this template to lowercase the $hookname in the URL
const RepoPage = `
{{ define "content" }}
	<br/><br/>
	<div>
		<div><b>{{.FullName}}</b></div>
		<div><b>Available validators</b>:<br>
		{{ range $hookname, $hook := .Hooks }}
			{{ $hookname }}
			{{ if eq $hook.State 0 }}
				[Enabled] <a href="/repos/{{$.FullName}}/{{ $hook.ID }}/disable">disable</a>
			{{ else }}
					[Disabled] <a href="/repos/{{$.FullName}}/{{ $hookname }}/enable">enable</a>
			{{ end }}
				<br>
			{{ end }}
		</div>
		<div>{{.Description}} {{.Website}}</div>
		<hr>
	</div>
{{ end }}
`
