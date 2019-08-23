package templates

// GenericResults is a template that requires only a header text, a badge, and
// content. The content is displayed in a <pre> block.
const GenericResults = `
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
		<pre>
			{{.Content}}
		</pre>
	</div>
{{end}}
`
