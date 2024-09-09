import React, { useState } from 'react';
import style from './Form.module.css';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function Form() {
    // Foydalanuvchi tomonidan kiritilgan ma'lumotlarni saqlash uchun state yaratamiz
    const [formData, setFormData] = useState({
        name: '',
        tel: '',
        email: '',
        link: ''
    });

    // Toast orqali xabar ko'rsatish funksiyalari
    const notifySuccess = () => toast.success("Ma'lumotlar yuborildi!");
    const notifyError = (message) => toast.error(`Xato: ${message}`);

    // Formadagi inputlarni boshqarish funksiyasi
    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    // Formani yuborish funksiyasi
    const handleSubmit = async (e) => {
        e.preventDefault();
    
        try {
            const response = await fetch('http://195.2.84.169:2005/partner', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    full_name: formData.name,
                    phone: formData.tel,
                    email: formData.email,
                    video_link: formData.link,
                }),
            });
    
            if (!response.ok) {
                const errorData = await response.json();
                notifyError(errorData.message || "Xato yuz berdi");
            } else {
                const successData = await response.json();
                notifySuccess(successData.message || "Ma'lumotlar muvaffaqiyatli yuborildi!");
            }
        } catch (error) {
            notifyError("Server bilan bog'lanishda muammo");
        }
    };

    return (
        <div className={style.wrap}>
            <div className='container'>
                <form onSubmit={handleSubmit} className={style.form}>
                    <h2 style={{ fontWeight: 'bold' }}>Videoroliklar tanlovida ishtirok eting</h2>
                    <p style={{ fontSize: 14, margin: '15px 0' }}>
                        <span style={{ color: 'red', marginRight: 3 }}>*</span>
                        Required
                    </p>
                    <h3 style={{ color: '#BA3B02', fontWeight: 'bold' }}>Ma'lumotlaringizni to‘ldiring</h3>
                    <label className={style.label}>
                        <p>1. Ismingiz <span style={{ color: 'red', marginLeft: 3 }}>*</span></p>
                        <input className={style.input} type="text" placeholder='Ismingizni kiriting' name="name" required value={formData.name} onChange={handleChange} />
                    </label>
                    <label className={style.label}>
                        <p>2. Telefon raqamingiz <span style={{ color: 'red', marginLeft: 3 }}>*</span></p>
                        <input className={style.input} type="tel" placeholder='Telefon raqamingizni kiriting' name="tel" required value={formData.tel} onChange={handleChange} />
                    </label>
                    <label className={style.label}>
                        <p>3. E-mail <span style={{ color: 'red', marginLeft: 3 }}>*</span></p>
                        <input className={style.input} type="email" placeholder='E-mailingizni kiriting' name="email" required value={formData.email} onChange={handleChange} />
                    </label>
                    <label className={style.label}>
                        <p>4. Video manzilni kiriting <span style={{ color: 'red', marginLeft: 3 }}>*</span></p>
                        <input className={style.input} type="text" placeholder='Video linkni kiriting' name="link" required value={formData.link} onChange={handleChange} />
                    </label>

                    <button type="submit" className={style.button}>Yuborish</button>
                    <ToastContainer />
                </form>
            </div>
        </div>
    );
}

export default Form;
