import '../estilos/ReservationCard.css';

const ReservationCard = ({ reservation, onCancel }) => {
  const formatDate = (date) => {
    return new Date(date).toLocaleDateString('es-ES');
  };

  const getStatusClass = (status) => {
    if (status === 'confirmed') return 'status-confirmed';
    if (status === 'pending') return 'status-pending';
    return 'status-cancelled';
  };

  const getStatusText = (status) => {
    if (status === 'confirmed') return 'Confirmada';
    if (status === 'pending') return 'Pendiente';
    return 'Cancelada';
  };

  return (
    <div className="reservation-card">
      <div className="reservation-header">
        <div className="reservation-info">
          <h3>{reservation.room?.name || 'Habitación'}</h3>
          <p><strong>Check-in:</strong> {formatDate(reservation.checkIn)}</p>
          <p><strong>Check-out:</strong> {formatDate(reservation.checkOut)}</p>
          <p>
            <strong>Estado:</strong>
            <span className={`status-badge ${getStatusClass(reservation.status)}`}>
              {getStatusText(reservation.status)}
            </span>
          </p>
          <p className="reservation-total"><strong>Total:</strong> ${reservation.total || reservation.room?.price}</p>
        </div>
        <button 
          onClick={() => onCancel(reservation.id)}
          className="btn-cancel"
        >
          Cancelar Reserva
        </button>
      </div>
    </div>
  );
};

export default ReservationCard;