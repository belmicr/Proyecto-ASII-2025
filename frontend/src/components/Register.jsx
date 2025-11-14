import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authService } from "../services/api";
import "../estilos/Register.css";

const Register = () => {
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    email: "",
    password: "",
    name: "",
  });

  const [error, setError] = useState("");

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const data = await authService.register(formData);

      
      if (data.error) {
        setError(data.error);
        return;
      }

      alert("Usuario registrado con éxito");
      navigate("/login");
    } catch (err) {
      console.error(err);
      setError("Error al registrar el usuario");
    }
  };

  return (
    <div className="register-container">
      <h2>Crear cuenta</h2>

      {error && <p className="error">{error}</p>}

      <form onSubmit={handleSubmit} className="register-form">
        <label>
          Nombre
          <input
            type="text"
            name="name"
            placeholder="Tu nombre"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </label>

        <label>
          Email
          <input
            type="email"
            name="email"
            placeholder="tuemail@ejemplo.com"
            value={formData.email}
            onChange={handleChange}
            required
          />
        </label>

        <label>
          Contraseña
          <input
            type="password"
            name="password"
            placeholder="********"
            value={formData.password}
            onChange={handleChange}
            required
          />
        </label>

        <button type="submit" className="btn-primary">
          Registrarme
        </button>
      </form>

      <p className="link">
        ¿Ya tenés cuenta? <a href="/login">Iniciar sesión</a>
      </p>
    </div>
  );
};

export default Register;
