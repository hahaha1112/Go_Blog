// 文档加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    console.log('GoBlog 已加载');

    // 添加当前年份到页脚
    const currentYearElement = document.querySelector('footer .container p');
    if (currentYearElement) {
        const year = new Date().getFullYear();
        currentYearElement.innerHTML = currentYearElement.innerHTML.replace('{{ .CurrentYear }}', year);
    }

    // 对表单进行简单验证
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function(event) {
            let isValid = true;
            
            // 检查所有必填字段
            const requiredFields = form.querySelectorAll('[required]');
            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    isValid = false;
                    field.classList.add('error');
                } else {
                    field.classList.remove('error');
                }
            });
            
            // 如果是注册表单，检查密码是否匹配
            if (form.action.includes('/register/process')) {
                const password = form.querySelector('#password');
                const confirmPassword = form.querySelector('#confirm_password');
                
                if (password && confirmPassword && password.value !== confirmPassword.value) {
                    isValid = false;
                    confirmPassword.classList.add('error');
                    alert('两次密码输入不一致');
                }
            }
            
            if (!isValid) {
                event.preventDefault();
            }
        });
    });

    // 添加响应式导航菜单
    const navToggle = document.querySelector('.nav-toggle');
    if (navToggle) {
        navToggle.addEventListener('click', function() {
            const nav = document.querySelector('nav ul');
            nav.classList.toggle('show');
        });
    }
}); 