package templates

// NotValidatedYet is a template that requires a header text, a badge, a
// content and two links. The content is displayed in a <pre> block.
const NotValidatedYet = `
{{define "content"}}
	<div class="repository file list">
		<div class="header-wrapper">
			<div class="ui container">
				<div class="ui vertically padded grid head">
					<div class="column">
						<div class="ui header">
							<div class="ui huge breadcrumb">
								<i class="mega-octicon octicon-repo"></i>
								{{.Header}}
								{{.Badge}}
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="ui tabs container">
			</div>
			<div class="ui tabs divider"></div>
		</div>
		<div class="ui container">
			<hr>
			<div>
				<pre>{{.Content}}</pre>
				<a href="{{.HrefURL1}}" alt="{{.HrefAlt1}}">{{.HrefText1}}</a> | 
				<a href="{{.HrefURL2}}" alt="{{.HrefAlt2}}">{{.HrefText2}}</a>
			</div>
		</div>
	</div>
{{end}}
`
