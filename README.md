# vpub-plus

Forum written in Go, fork of [vpub](https://sr.ht/~m15o/vpub/) with a lot of small additions.

## Examples

<table>
  <tr>
    <td>
      <a href="https://github.com/user-attachments/assets/4ded261a-f2c3-4f8f-b474-be268aa61ff7">
        <img alt="Status Cafe Forum - vpub instance" width="640" src="https://github.com/user-attachments/assets/4ded261a-f2c3-4f8f-b474-be268aa61ff7">
      </a>
      <p><a href="https://forum.status.cafe">Status Cafe Forum - vpub instance</a></p>
    </td>
    <td>
      <a href="https://github.com/user-attachments/assets/6e65c795-6ef6-40f3-b222-18f4c3f48548">
        <img alt="Vpub Plus Forum - vpub-plus instance" width="640" src="https://github.com/user-attachments/assets/6e65c795-6ef6-40f3-b222-18f4c3f48548">
      </a>
      <p><a href="https://vpub.mysh.dev">Vpub Plus Forum - vpub-plus instance</a></p>
    </td>
  </tr>
</table>

## Installation

See installation instructions on the forum itself: https://vpub.mysh.dev/topics/3

If you want to run it from binary and with Supabase: https://vpub.mysh.dev/topics/5

Looking for installation instructions of original vpub? You can use this tutorial: https://vpub.mysh.dev/topics/7

## Screenshots and Features

| Feature                                                                                                                                                     | GIF                                                                                                                                                  |
| :---------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| **CSS Customization** – customzie forum appearance to your liking: increase default width, customize colors with variables and such                         | ![GIF shows a demo of customization of CSS in admin panel](https://github.com/user-attachments/assets/a6085b19-86bd-4530-9910-8582de5dc830)          |
| **Footer customization** – add whatever HTML you want into the footer section. For example, a preload of links using HTMX                                   | ![GIF shows a demo of customization of a footer tag in admin panel](https://github.com/user-attachments/assets/a4444696-1896-4db6-b3c8-97006222f5f3) |
| **Changing of render engines** – you can change render engines on-the-fly. Out of the box vpub-plus supports `vanilla` and `blackfriday` markdown renderers | ![GIF shows a demo of changing of a renderer in admin panel](https://github.com/user-attachments/assets/a68918d2-8c4a-408a-882b-de2557e988b9)        |
| **Search** – that is something that everyone likes to have embedded in their software                                                                       | ![GIF shows a demo of search functionality](https://github.com/user-attachments/assets/eb8ee1c1-f5a9-4094-8e04-6817d2be1e64)                         |


### Admin Panel Screenshots

It's pretty minimal, but it has all the needed features.

| Panel                                                                                                                                                                             | Screenshot                                                                                                 |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| **General forum settings** – you can change name, absolute url, site language, provide custom CSS and Footer, change pagination limit and change rendering engine for forum posts | ![General forum settings](https://github.com/user-attachments/assets/de956518-de47-48eb-87dc-171917923cfe) |
| **Keys management** – create and delete keys (used for registration)                                                                                                              | ![Keys management](https://github.com/user-attachments/assets/83189a9b-0170-45fd-a7b0-906ae1a3e785)        |
| **Boards management** – add, edit, delete boards for forums. Change their position, lock status and their content                                                                 | ![Boards management](https://github.com/user-attachments/assets/b932c9fd-3303-4472-a9cb-830b08f03abb)      |
| **Forums management** – add, edit, delete forums. Change position and lock status                                                                                                 | ![Forums management](https://github.com/user-attachments/assets/d1677e7d-81d6-4c9a-9bb5-09543f2d90a9)      |
| **Users management** – edit and delete users. Change their name, picture and about section                                                                                        | ![Users management](https://github.com/user-attachments/assets/d7903f04-e7d4-4b4c-9edc-a11dd8c328d6)       |
