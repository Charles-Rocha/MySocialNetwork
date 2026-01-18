$('#formulario-cadastro').on('submit', criarUsuario);

function criarUsuario(evento){
    evento.preventDefault();
    //console.log("Dentro da função CriarUsuario");

    if ($('#senha').val() != $('#confirmar-senha').val()){
        alert("As senhas são diferentes!");
        return;
    }

    $.ajax({
        url: "/usuarios",
        method: "POST",
        data: {
            nome: $('#nome').val(),
            email: $('#email').val(),
            nick: $('#nick').val(),
            senha: $('#senha').val()
        }
    }).done(function(){//O Ajax sabe identificar se deu certo ou não baseado no StatusCode. 200 201 204
        alert('Usuário cadastrado com sucesso!');
    }).fail(function(erro){
        console.log(erro);
        alert('Erro ao tentar cadastrar usuário!');
    });
}