$('#formulario-cadastro').on('submit', criarUsuario);

function criarUsuario(evento){
    evento.preventDefault();
    //console.log("Dentro da função CriarUsuario");

    if ($('#senha').val() != $('#confirmar-senha').val()){
        Swal.fire("Atenção", "As senhas informadas não são iguais.", "warning");        
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
        Swal.fire("Sucesso", "Usuário cadastrado com sucesso.", "success")
            .then(function(){
                $.ajax({
                    url: "/login",
                    method: "POST",
                    data: {
                        email: $('#email').val(),
                        senha: $('#senha').val()
                    }
                }).done(function(){
                    window.location = "/home";
                }).fail(function(){
                    Swal.fire("Atenção", "Erro ao tentar autenticar o usuário.", "error");    
                })
            })         
    }).fail(function(erro){
        console.log(erro);
        Swal.fire("Atenção", "Erro ao tentar cadastrar usuário.", "error");        
    });
}