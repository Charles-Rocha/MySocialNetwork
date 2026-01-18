$('#login').on('submit', fazerLogin);

function fazerLogin(evento){
    alert('Dentro da função Login');
    evento.preventDefault();
    console.log("Dentro da função CriarUsuario");
    
    $.ajax({
        url: "/login",
        method: "POST",
        data: {
            email: $('#email').val(),
            senha: $('#senha').val(),
        }
    }).done(function(){//O Ajax sabe identificar se deu certo ou não baseado no StatusCode. 200 201 204
        alert('Login OK')
        window.location = "/home"
    }).fail(function(erro){  
        console.log(erro)      
        alert('Usuário ou senha inválidos!')
    });
}

