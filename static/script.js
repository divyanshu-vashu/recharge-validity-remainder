document.getElementById('simForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const rechargeDate = new Date(document.getElementById('rechargeDate').value);
    
    const simData = {
        name: document.getElementById('simName').value,
        number: document.getElementById('simNumber').value,
        lastRechargeDate: rechargeDate.toISOString()
    };

    try {
        const response = await fetch('/api/sims', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(simData)
        });
        
        if (response.ok) {
            const result = await response.json();
            console.log('Added SIM:', result);  // Debug log
            loadSims();
            document.getElementById('simForm').reset();
        } else {
            const error = await response.text();
            console.error('Error adding SIM:', error);  // Debug log
            alert('Failed to add SIM: ' + error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add SIM');
    }
});

document.getElementById('syncButton').addEventListener('click', () => {
    loadSims();
});

document.getElementById('filterType').addEventListener('change', () => {
    loadSims();
});

async function loadSims() {
    try {
        const response = await fetch('/api/sims');
        const sims = await response.json();
        
        const simList = document.getElementById('simList');
        simList.innerHTML = '';
        
        sims.forEach(sim => {
            const card = document.createElement('div');
            card.className = 'card mb-3';
            card.innerHTML = `
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start">
                        <div>
                            <h5 class="card-title">${sim.name} - ${sim.number}</h5>
                            <p class="mb-2">
                                Last Recharge: ${new Date(sim.lastRechargeDate).toLocaleDateString()}
                                <button class="btn btn-sm btn-outline-primary ms-2" onclick="openEditModal(${sim.id}, '${sim.lastRechargeDate}')">
                                    Edit
                                </button>
                            </p>
                            <p>Recharge Validity: ${new Date(sim.rechargeValidity).toLocaleDateString()}</p>
                            <p>Incoming Call Validity: ${new Date(sim.incomingValidity).toLocaleDateString()}</p>
                            <p>SIM Expiry: ${new Date(sim.simExpiry).toLocaleDateString()}</p>
                        </div>
                    </div>
                </div>
            `;
            simList.appendChild(card);
        });
    } catch (error) {
        console.error('Error:', error);
    }
}

// Add these new functions for editing
function openEditModal(simId, currentDate) {
    document.getElementById('editSimId').value = simId;
    document.getElementById('editRechargeDate').value = new Date(currentDate).toISOString().split('T')[0];
    new bootstrap.Modal(document.getElementById('editModal')).show();
}

document.getElementById('saveEdit').addEventListener('click', async () => {
    const simId = document.getElementById('editSimId').value;
    const newDate = document.getElementById('editRechargeDate').value;

    try {
        const response = await fetch(`/api/sims/${simId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                lastRechargeDate: newDate
            })
        });

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('editModal')).hide();
            loadSims();
        } else {
            alert('Failed to update recharge date');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update recharge date');
    }
});

function getStatus(sim, currentDate) {
    const rechargeDate = new Date(sim.rechargeValidity);
    const daysUntilExpiry = Math.ceil((rechargeDate - currentDate) / (1000 * 60 * 60 * 24));
    
    if (daysUntilExpiry < 0) return 'Expired';
    if (daysUntilExpiry <= 3) return 'Expiring Soon';
    return 'Active';
}

function getStatusBadgeClass(sim, currentDate) {
    const rechargeDate = new Date(sim.rechargeValidity);
    const daysUntilExpiry = Math.ceil((rechargeDate - currentDate) / (1000 * 60 * 60 * 24));
    
    if (daysUntilExpiry < 0) return 'bg-danger';
    if (daysUntilExpiry <= 3) return 'bg-warning';
    return 'bg-success';
}

// Load sims when page loads
loadSims();