document.addEventListener('DOMContentLoaded', function() {
    //attempt email obfuscation
    setTimeout(function() {
        var emailLink = document.getElementById("email");
        if (emailLink) {
        emailLink.href = atob("bWFpbHRvOnNpbW9uQHNhd2F0a2lucy5jb20=");
        }
    }, 3000);
});