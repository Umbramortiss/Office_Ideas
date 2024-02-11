// Example JavaScript code to handle form submission
document.addEventListener('DOMContentLoaded', function (){
  const userForm = document.getElementById('userForm');
userForm.addEventListener('submit',async(event) => {

    event.preventDefault()
    const formData = new FormData(userForm);
    const formDataObject = {
        Fname: 'John',
        Lname: 'Doe',
        Email: 'john.doe@example.com',
        Sugg: 'This is a suggestion.',
    };

    formData.forEach((value, key) => {
        formDataObject[key] = value;
    });


  // Send the data to the backend using Fetch API
  try {

    const response = await fetch('/api/submit', {

      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formDataObject),
    });

    if (response.ok) {
      // Handle successful response from the backend
      console.log('User created successfully!');
    } else {
      // Handle errors from the backend
      console.error('Error creating user:', response.statusText);
    }
  } catch (error) {
    // Handle network errors
    console.error('Network error:', error);
  }
});
});
