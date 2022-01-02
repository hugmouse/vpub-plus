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
	"create_post": `{{ define "title" }}New Post{{ end }}

{{ define "content" }}
    <h2>New Post</h2>
    <form action="/posts/save" method="post">
        {{ .csrfField }}
        {{ template "post_form" .form }}
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

{{ .motd }}

<nav class="actions">
    <p>
        {{ if .logged }}
        <a href="/posts/new">write</a>
        {{ end }}
        <a href="/feed.atom">follow</a>
    </p>
</nav>

{{ template "topics" . }}

<section>
{{ template "posts" .posts }}
{{ if .hasMore }}
<a href="/page/2">More</a>
{{ end }}
</section>

<section>
    <h2>Users</h2>
{{ range .users }}<a href="/~{{ . }}">{{ . }}</a> {{ end }}
</section>
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
	"post": `{{ define "content"}}
<h1>{{ .post.Title }}</h1>
<div class="meta">
    {{ with .post }}
    <ul class="key-value">
        <li><span class="key">From: </span><span class="value"><a href="/~{{ .User }}">{{ .User }}</a></span></li>
        <li><span class="key">On: </span><span class="value">{{ .Date }}</span></li>
        {{ if .Topic }}<li><span class="key">Topic: </span><span class="value"><a href="/topics/{{ .Topic }}">{{ .Topic }}</a></span></li>{{ end }}
    </ul>
    {{ end }}
</div>
<div class="content">{{ gmi2html .content }}</div>
{{- if eq .logged .post.User }}
<p>
    <a href="/posts/{{ .post.Id }}/edit">Edit</a>
    <a href="/posts/{{ .post.Id }}/remove">Remove</a>
</p>
{{- end }}
{{ if .logged }}
<form action="/posts/{{ .post.Id }}/reply" method="post">
    {{ .csrfField }}
    <textarea name="reply"></textarea>
    <input type="submit" value="Reply">
</form>
{{ end }}
{{ template "reply" .replies }}
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
            <input type="text" id="name" name="name" autocomplete="off" value="{{ .form.Username }}"/>
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
                <li><span class="key">From: </span><span class="value"><a href="/~{{ .post.User }}">{{ .post.User }}</a></span></li>
                <li><span class="key">On: </span><span class="value">{{ .post.Date }}</span></li>
                <li><span class="key">Post: </span><span class="value"><a href="/posts/{{ .post.Id }}">{{ .post.Title }}</a></span></li>
                <li><span class="key">Parent: </span><span class="value">
                    {{ if .reply.ParentId }}<a href="/replies/{{ .reply.ParentId }}">view</a>{{ else }}view{{end}}
                </span></li>
            </ul>
        </div>
        {{ gmi2html .reply.Content }}
    </article>
    <section>
        <form action="/replies/{{ .reply.Id }}/save" method="post">
            {{ .csrfField }}
            <textarea name="reply" autofocus></textarea>
            <input type="submit" value="Submit">
        </form>
        {{ template "reply" .reply.Thread }}
    </section>
{{ end }}
`,
	"topic": `{{ define "content" }}
<h1>{{ .topic }}</h1>

<nav class="actions">
    <p>
        {{ if .logged }}
        <a href="/posts/new?topic={{ .topic }}">write</a>
        {{ end }}
        <a href="/topics/{{ .topic }}/feed.atom">follow</a>
    </p>
</nav>

{{ template "topics" . }}

<section class="posts">
    {{ template "posts" .posts }}
    {{ if .hasMore }}
    <a href="/page/2?topic={{ .topic }}">More</a>
    {{ end }}
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
