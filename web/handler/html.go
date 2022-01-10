// Code generated by go generate; DO NOT EDIT.

package handler

var TplMap = map[string]string{
	"account": `{{ define "title" }}Account{{ end }}

{{ define "content" }}
<h1>Account</h1>

<form action="/save-about" method="post">
    {{ .csrfField }}
    <div class="field">
        <label for="about">About</label>
        <textarea name="about" id="about" autofocus>{{ .about }}</textarea>
    </div>
    <input type="submit" value="Submit">
</form>
</section>
{{ end }}
`,
	"admin_board": `{{ define "content"}}
<h1>Admin - Boards</h1>
<p><a href="/admin/boards/new">New board</a></p>
<table>
    <thead>
    <tr>
        <th class="grow">Board</th>
        <th>Edit</th>
    </tr>
    </thead>
    <tbody>
    {{ range .boards }}
    <tr>
        <td colspan="grow">
            <a href="/boards/{{ .Id }}">{{ .Name }}</a><br>{{ .Description }}
        </td>
        <td class="center"><a href="/admin/boards/{{ .Id }}/edit">Edit</a></td>
    </tr>
    {{ end }}
    </tbody>
</table>
{{ end }}`,
	"admin_board_create": `{{ define "title" }}New board{{ end }}
{{ define "content" }}
<h2>Admin - Create board</h2>
<form action="/admin/boards/save" method="post">
  {{ .csrfField }}
  {{ template "board_form" .form }}
  <input type="submit" value="Submit">
</form>
{{ end }}
`,
	"admin_board_edit": `{{ define "title" }}New board{{ end }}
{{ define "content" }}
<h2>Admin - Edit board</h2>
<form action="/admin/boards/{{ .board.Id }}/update" method="post">
    {{ .csrfField }}
    {{ template "board_form" .form }}
    <input type="submit" value="Submit">
</form>
{{ end }}
`,
	"board": `{{ define "breadcrumb" }} > {{ .board.Name }}{{ end }}
{{ define "content" }}
<h1>{{ .board.Name }}</h1>

<!--<p>{{ .board.Description }}</p>-->

<nav class="actions">
    <p>
        {{ if .logged }}
        <a href="/boards/{{ .board.Id }}/new-topic">new topic</a>
        {{ end }}
<!--        <a href="TODO">follow</a>-->
    </p>
</nav>

<section>
    <table>
        <thead>
        <tr>
            <th class="grow">Subject</th>
            <th>Author</th>
            <th>Replies</th>
            <th>Updated</th>
        </tr>
        </thead>
        <tbody>
        {{ range .topics }}
        <tr>
            <td colspan="grow"><a href="/topics/{{ .Id }}">{{ .Subject }}</a></td>
            <td class="center"><a href="/~{{ .Author }}">{{ .Author }}</a></td>
            <td class="center">{{ .Replies }}</td>
            <td class="center">{{ iso8601 .UpdatedAt }}</td>
        </tr>
        {{ end }}
        </tbody>
    </table>
</section>
{{ end }}`,
	"confirm_remove_post": `{{ define "content" }}
    Are you sure you you want to delete the following post?
    <p>{{ gmi2html .post.Content }}</p>
    <form action="/posts/{{ .post.Id }}/remove" method="post">
        {{ .csrfField }}
        <input type="submit" value="Submit">
    </form>
{{ end }}`,
	"confirm_remove_reply": `{{ define "content" }}
    Are you sure you you want to delete the following reply?
    <p>{{ gmi2html .reply.Content }}</p>
    <form action="/replies/{{ .reply.Id }}/remove" method="post">
        {{ .csrfField }}
        <input type="submit" value="Submit">
    </form>
{{ end }}`,
	"create_post": `{{ define "title" }}New Thread{{ end }}
{{ define "breadcrumb" }} > <a href="/topics/{{ .form.Topic.Id }}">{{ .form.Topic.Name }}</a>{{ end }}
{{ define "content" }}
    <h2>New Thread</h2>
    <form action="/posts/save" method="post">
        {{ .csrfField }}
        {{ template "post_form" .form }}
        <input type="submit" value="Submit">
    </form>
{{ end }}
`,
	"create_topic": `{{ define "title" }}New Thread{{ end }}
{{ define "breadcrumb" }} > <a href="/boards/{{ .board.Id }}">{{ .board.Name }}</a>{{ end }}
{{ define "content" }}
<h2>New topic</h2>
<form action="/boards/{{ .board.Id }}/save-topic" method="post">
    {{ .csrfField }}
    <input type="hidden" name="boardId" value="{{ .board.Id }}">
    {{ template "post_form" .form.PostForm }}
    <input type="submit" value="Submit">
</form>
{{ end }}
`,
	"edit_post": `{{ define "title" }}Edit Post{{ end }}

{{ define "content" }}
    <h2>Edit Post</h2>
    <form action="/posts/{{ .post.Id }}/update" method="post">
        {{ .csrfField }}
        {{ template "post_form" .form }}
        <input type="submit" value="Update">
    </form>
{{ end }}
`,
	"edit_reply": `{{ define "content" }}
    <form action="/replies/{{ .reply.Id }}/update" method="post" class="form">
        {{ .csrfField }}
        <label for="reply">reply</label><textarea name="reply" id="reply">{{ .form.Content }}</textarea>
        <input type="submit" value="Submit">
    </form>
{{ end }}
`,
	"index": `{{ define "content"}}
<h1>{{ .boardTitle }}</h1>

<!--<nav class="actions">-->
<!--    <p>-->
<!--        <a href="/feed.atom">follow</a>-->
<!--    </p>-->
<!--</nav>-->

<table>
    <thead>
        <tr>
            <th class="grow">Board</th>
            <th>Topics</th>
            <th>Posts</th>
            <th>Updated</th>
        </tr>
    </thead>
    <tbody>
        {{ range .boards }}
        <tr>
            <td colspan="grow">
                <a href="/boards/{{ .Id }}">{{ .Name }}</a><br>{{ .Description }}
            </td>
            <td class="center">{{ .Topics }}</td>
            <td class="center">{{ .Posts }}</td>
            <td class="center">{{ iso8601 .UpdatedAt }}</td>
        </tr>
        {{ end }}
    </tbody>
</table>
{{ end }}`,
	"login": `{{ define "title" }}Login{{ end }}

{{ define "content" }}
<h2>Login</h2>
<form action="/login" method="post" class="auth-form">
    {{ .csrfField }}
    <div class="field">
        <label for="name">Username</label>
        <input type="text" name="name" id="name" autocomplete="off" required/>
    </div>
    <div class="field">
        <label for="password">Password</label>
        <input type="password" name="password" id="password" required/>
    </div>
    <input type="submit" value="Login">
</form>
{{ end }}`,
	"notifications": `{{ define "content" }}
<h1>New replies</h1>
<p><a href="/notifications/mark-all-read">mark all as read</a></p>
{{ range .notifications }}
<div>
    <div class="meta">
        <ul class="key-value">
            <li><span class="key">From: </span><span class="value"><a href="/~{{ .Reply.User }}">{{ .Reply.User }}</a></span></li>
            <li><span class="key">On: </span><span class="value">{{ .Reply.Date }}</span></li>
            <li><span class="key">Post: </span><span class="value"><a href="/posts/{{ .Reply.PostId }}">{{ .Reply.PostTitle }}</a></span></li>
            <li><span class="key">Parent: </span><span class="value">
                        {{ if .Reply.ParentId }}<a href="/replies/{{ .Reply.ParentId }}">view</a>{{ else }}view{{end}}
            </span></li>
        </ul>
    </div>
    <div class="content">{{ gmi2html .Reply.Content }}</div>
    <p>
        <a href="/replies/{{ .Reply.Id }}">reply</a>
        <a href="/notifications/{{ .Id }}/mark-read">mark as read</a>
    </p>
</div>
{{ end }}
{{ end }}`,
	"paginate": `{{ define "content" }}
<p>Page {{ .page }}{{ if .topic }} of <a href="/topics/{{ .topic }}">{{ .topic }}</a>{{ end }}</p>
<section>
    {{ template "posts" .posts }}

    {{ if .hasMore }}
    {{ if .topic }}
    <a href="/page/{{ .nextPage }}?topic={{ .topic }}">More</a>
    {{ else }}
    <a href="/page/{{ .nextPage }}">More</a>
    {{ end }}
    {{ end }}
</section>
{{ end }}`,
	"post": `{{ define "breadcrumb" }} > <a href="/topics/{{ .post.Topic }}">{{ .post.Topic }}</a>{{ end }}
{{ define "content"}}
<!--<table class="thread">-->
<!--    <tr>-->
<!--        <td>{{ .post.User }}</td>-->
<!--        <td>{{ gmi2html .content }}</td>-->
<!--    </tr>-->
<!--    {{ range .replies }}-->
<!--    <tr>-->
<!--        <td>{{ .User }}</td>-->
<!--        <td>{{ gmi2html .Content }}</td>-->
<!--    </tr>-->
<!--    {{ end }}-->
<!--</table>-->
<!--<form action="/posts/{{ .post.Id }}/reply" method="post">-->
<!--    {{ .csrfField }}-->
<!--    <div class="field">-->
<!--        <textarea name="reply"></textarea>-->
<!--    </div>-->
<!--    <input type="submit" value="Reply">-->
<!--</form>-->
<h1>{{ .post.Title }}</h1>
<!--<ol class="thread">-->
<!--    <li class="post">-->
<!--        <table>-->
<!--            <tr>-->
<!--                <td class="post-aside">{{ .post.User }}</td>-->
<!--                <td class="post-content">{{ gmi2html .content }}</td>-->
<!--            </tr>-->
<!--        </table>-->
<!--    </li>-->
<!--    {{ range .replies }}-->
<!--    <li class="post">-->
<!--        <table>-->
<!--            <tr>-->
<!--                <td class="post-aside">{{ .User }}</td>-->
<!--                <td class="post-content">{{ gmi2html .Content }}</td>-->
<!--            </tr>-->
<!--        </table>-->
<!--    </li>-->
<!--    {{ end }}-->
<!--</ol>-->
<table class="thread">
    <tr class="post">
        <td class="post-aside">
            <p>{{ .post.User }}</p>
            <p>{{ timeAgo .post.CreatedAt }}</p>
        </td>
        <td class="post-content">
            {{ gmi2html .content }}
        </td>
    </tr>
    {{ range .replies }}
    <tr class="post">
        <td class="post-aside">
            <p>{{ .User }}</p>
            <p>{{ timeAgo .CreatedAt }}</p>
        </td>
        <td class="post-content">
            {{ gmi2html .Content }}
        </td>
    </tr>
    {{ end }}
</table>
<form action="/posts/{{ .post.Id }}/reply" method="post">
    {{ .csrfField }}
    <div class="field">
        <textarea name="reply"></textarea>
    </div>
    <input type="submit" value="Reply">
</form>
<!--<h1>{{ .post.Subject }}</h1>-->
<!--<div class="meta">-->
<!--    {{ with .post }}-->
<!--    <ul class="key-value">-->
<!--        <li><span class="key">From: </span><span class="value"><a href="/~{{ .User }}">{{ .User }}</a></span></li>-->
<!--        <li><span class="key">On: </span><span class="value">{{ timeAgo .CreatedAt }} ({{ .Date }})</span></li>-->
<!--        {{ if .Topic }}<li><span class="key">Topic: </span><span class="value"><a href="/topics/{{ .Topic }}">{{ .Topic }}</a></span></li>{{ end }}-->
<!--    </ul>-->
<!--    {{ end }}-->
<!--</div>-->
<!--<div class="content">{{ gmi2html .content }}</div>-->
<!--{{- if eq .logged .post.User }}-->
<!--<p>-->
<!--    <a href="/posts/{{ .post.Id }}/edit">Edit</a>-->
<!--    <a href="/posts/{{ .post.Id }}/remove">Remove</a>-->
<!--</p>-->
<!--{{- end }}-->
<!--{{ if .logged }}-->
<!--<form action="/posts/{{ .post.Id }}/reply" method="post">-->
<!--    {{ .csrfField }}-->
<!--    <div class="field">-->
<!--        <textarea name="reply"></textarea>-->
<!--    </div>-->
<!--    <input type="submit" value="Reply">-->
<!--</form>-->
<!--{{ end }}-->
<!--{{ template "reply" .replies }}-->
{{ end }}`,
	"register": `{{ define "title" }}Register{{ end }}

{{ define "content" }}
    <h2>Register</h2>
    {{ if .error }}
    <p class="error">{{ .error }}</p>
    {{ end }}
    <form action="/register" method="post" class="auth-form">
        {{ .csrfField }}
        <div class="field">
            <label for="name">Username</label>
            <input type="text" id="name" name="name" autocomplete="off" value="{{ .form.Username }}" maxlength="15"/>
        </div>
        <div class="field">
            <label for="password">Password</label>
            <input type="password" id="password" name="password"/>
        </div>
        <div class="field">
            <label for="confirm">Confirm password</label>
            <input type="password" id="confirm" name="confirm" required/>
        </div>
<!--        <div class="field">-->
<!--            <label for="key">Key</label>-->
<!--            <input type="text" id="key" name="key"/>-->
<!--        </div>-->
        <input type="submit" value="Submit">
    </form>
{{ end }}
`,
	"reply": `{{ define "content" }}
    <h1>Reply</h1>
    <article>
        <div class="meta">
            <ul class="key-value">
                <li><span class="key">From: </span><span class="value"><a href="/~{{ .reply.User }}">{{ .reply.User }}</a></span></li>
                <li><span class="key">On: </span><span class="value">{{ .post.Date }}</span></li>
                <li><span class="key">Post: </span><span class="value"><a href="/posts/{{ .post.Id }}">{{ .post.Title }}</a></span></li>
                <li><span class="key">Parent: </span><span class="value">
                    {{ if .reply.ParentId }}<a href="/replies/{{ .reply.ParentId }}">view</a>{{ else }}view{{end}}
                </span></li>
            </ul>
        </div>
        <div class="content">{{ gmi2html .reply.Content }}</div>
    </article>
    <section>
        <form action="/replies/{{ .reply.Id }}/save" method="post">
            {{ .csrfField }}
            <div class="field">
                <textarea name="reply" autofocus></textarea>
            </div>
            <input type="submit" value="Submit">
        </form>
        {{ template "reply" .reply.Thread }}
    </section>
{{ end }}
`,
	"topic": `{{ define "breadcrumb" }} > <a href="/boards/{{ .board.Id }}">{{ .board.Name }}</a>{{ end }}
{{ define "content"}}
<!--<h1>{{ .topic.Subject }}</h1>-->
<br>
<table>
    {{ range .posts }}
    <tr class="post">
        <td class="post-aside">
            <p>{{ .User }}</p>
            <p>{{ timeAgo .CreatedAt }}</p>
        </td>
        <td class="post-content">
            {{ if eq $.topic.FirstPostId .Id }}<h1>{{ .Title }}</h1>{{ end }}
            {{ gmi2html .Content }}
            {{ if hasPermission .User }}
            <p><a href="/posts/{{ .Id }}/edit">edit</a> <a href="/posts/{{ .Id }}/remove">remove</a></p>
            {{ end }}
        </td>
    </tr>
    {{ end }}
</table>
<section style="margin-top: 1em;">
    <form action="/posts/save" method="post">
        {{ .csrfField }}
        <input type="hidden" name="topicId" value="{{ .topic.Id }}">
        <input type="hidden" name="subject" value="Re: {{ .topic.Subject }}">
        <div class="field">
            <label for="content">Reply to this topic</label>
            <textarea name="content" id="content" style="height: 100px;"></textarea>
        </div>
        <input type="submit" value="Reply">
    </form>
</section>
{{ end }}`,
	"user_posts": `{{ define "content" }}
<h1>{{ .user.Name }}</h1>
<div class="content">{{ gmi2html .user.About }}</div>
<section class="posts">
{{ template "posts" .posts }}

{{if .showMore }}
<a href="/~{{ .user.Name }}?page={{ .nextPage }}">More</a>
{{ end }}
</section>
{{ end }}
`,
}
