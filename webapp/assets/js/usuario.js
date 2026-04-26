$('#parar-de-seguir').on('click', pararDeSeguir);
$('#seguir').on('click', seguir);
$('#editar-usuario').on('submit', editar);
$('#atualizar-senha').on('submit', atualizarSenha);
$('#deletar-usuario').on('click', deletarUsuario);


function pararDeSeguir(){    
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled', true);

    $.ajax({
        url: `/usuarios/${usuarioId}/parar-de-seguir`,
        method: "POST",        
    }).done(function(){
        window.location = `/usuarios/${usuarioId}`;        
    }).fail(function(){
        Swal.fire("Atenção", "Erro ao parar de seguir usuário", "error");
        $('#parar-de-seguir').prop('disabled', false);
    });
}

function seguir(){    
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled', true);

    $.ajax({
        url: `/usuarios/${usuarioId}/seguir`,
        method: "POST",        
    }).done(function(){
        window.location = `/usuarios/${usuarioId}`;        
    }).fail(function(){
        Swal.fire("Atenção", "Erro ao seguir usuário", "error");
        $('#seguir').prop('disabled', false);
    });
}

function editar(evento){
    evento.preventDefault();

    $.ajax({
        url: "/editar-usuario",
        method: "PUT",
        data: {
            nome: $('#nome').val(),
            email: $('#email').val(),
            nick: $('#nick').val(),
        }
    }).done(function(){
        Swal.fire("Sucesso!", "Usuário atualizado com sucesso", "success")
            .then(function(){
                window.location = "/perfil";
            });
    }).fail(function(){
        Swal.fire("Atenção!", "Erro ao atualizar o usuário", "error")               
    })
}

function atualizarSenha(evento){
    evento.preventDefault();

    if ($('#nova-senha').val() != $('#confirmar-senha').val()){
        Swal.fire("Atenção!", "As senha informadas não coincidem", "warning");
        return;
    }

    $.ajax({
        url: "/atualizar-senha",
        method: "POST",
        data: {
            atual: $('#senha-atual').val(),
            nova: $('#nova-senha').val(),            
        }
    }).done(function(){
        Swal.fire("Sucesso!", "Senha atualizada com sucesso", "success")
            .then(function(){
                window.location = "/perfil";
            });
    }).fail(function(){
        Swal.fire("Atenção!", "Erro ao atualizar a senha do usuário", "error")               
    })
}

function deletarUsuario(evento){
    Swal.fire({
        title: "Atenção",
        text: "Tem certeza que deseja apagar a sua conta? Essa ação é irreversível!",
        showCancelButton: true,
        cancelButtonText: "Cancelar",
        icon: "warning"
    }).then(function(confirmacao){
        if (confirmacao.value){
            $.ajax({
                url: "/deletar-usuario",
                method: "DELETE",                
            }).done(function(){
                Swal.fire("Sucesso!", "Usuário excluído com sucesso", "success")
                    .then(function(){
                        window.location = "/logout";
                    });
            }).fail(function(){
                Swal.fire("Atenção!", "Erro ao excluir o usuário", "error")               
            })
        }
    })    
}