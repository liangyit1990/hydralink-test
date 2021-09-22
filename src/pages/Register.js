import React, { useState, useEffect } from 'react';


const Register = () => {
    
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const submit = async (e)=> {
        e.preventDefault();
        const response = await fetch('http://localhost:4000/api/v1/users/signup',{
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                firstName,
                lastName,
                email,
                password
            })
        });

        const content = await response.json();
        console.log(content);
    }    


    // const submit = async (e) => {
    //     e.preventDefault();
        
    //     const response = await fetch('http://localhost:4000/api/v1/users/signup',{
    //         method: 'POST',
    //         headers: {'Content-Type': 'application/json'},
    //         body: JSON.stringify(value:{
    //             name,
    //             email,
    //             password
    //         })
    //     })
        
    
   
    return (

        <form onSubmit={submit}>
          
          <h1 className="h3 mb-3 fw-normal">Please register</h1>
          
          <input className="form-control" placeholder="First Name" 
            onChange={e =>setFirstName(e.target.value)}/>
        <input className="form-control" id="floatingInput" placeholder="Last Name" 
            onChange={e =>setLastName(e.target.value)}/>  

          <div className="form-floating">
            <input type="email" className="form-control"  placeholder="name@example.com"
            onChange={e =>setEmail(e.target.value)}/>
            <label for="floatingInput">Email address</label>
          </div>
          <div className="form-floating">
            <input type="password" className="form-control" placeholder="Password"
            onChange={e =>setPassword(e.target.value)}/>
            <label for="floatingPassword">Password</label>
          </div>

          
          <button className="w-100 btn btn-lg btn-primary" type="submit">Register</button>
          
        </form>
    )
}

export default Register
