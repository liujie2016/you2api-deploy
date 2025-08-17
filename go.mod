module you2api

go 1.22.2



curl https://you2api-deploy-git-main-liujiers-projects.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IlNLMmpJbnU3SWpjMkp1eFJad1psWHBZRUpQQkFvIiwidHlwIjoiSldUIn0.eyJhbXIiOlsiZW1haWwiXSwiYXV0aDBJZCI6Imdvb2dsZS1vYXV0aDJ8MTAxMTMwNDgwMTAzMzgzMzcxNDIzIiwiY3JlYXRlVGltZSI6MTcyMjMwNzcxMywiZHJuIjoiRFMiLCJlbWFpbCI6ImxpdWppZTE5OTYzMTBAZ21haWwuY29tIiwiZXhwIjoxNzU2NjUyMDI1LCJnaXZlbk5hbWUiOiJKaWUiLCJpYXQiOjE3NTU0NDI0MjUsImlzcyI6IlAyakludHRSTXVYcHlZWk1iVmNzYzRDOVowUlQiLCJsYXN0TmFtZSI6IkxpdSIsImxvZ2luSWRzIjpbImxpdWppZTE5OTYzMTBAZ21haWwuY29tIl0sIm5hbWUiOiIiLCJwaWN0dXJlIjoiIiwicmV4cCI6IjIwMjYtMDgtMTZUMTQ6NTM6NDVaIiwic3R5dGNoSWQiOiJ1c2VyLWxpdmUtYmE4OWQxNzctZDgxZi00YzYzLWJlNjMtMWYzNTUwZjI0ZjlkIiwic3ViIjoiVTJqd290TTZQN0o0cmZJczhvZzdGZjBVaDZIUiIsInN1YnNjcmlwdGlvblRpZXIiOiIiLCJ0ZW5hbnRDdXN0b21BdHRyaWJ1dGVzIjp7ImlzRW50ZXJwcmlzZSI6Int7dGVuYW50LmN1c3RvbUF0dHJpYnV0ZXMuaXNFbnRlcnByaXNlfX0iLCJuYW1lIjoie3t0ZW5hbnQubmFtZX19In0sInRlbmFudEludml0YXRpb24iOiIiLCJ0ZW5hbnRJbnZpdGVyIjoiIiwidXNlcklkIjoiVTJqd290TTZQN0o0cmZJczhvZzdGZjBVaDZIUiIsInZlcmlmaWVkRW1haWwiOnRydWV9.JN2yaXeRMP9IOh19grD6MVdTgj81Vlj2tIbmw4PomeQ5yJB_Y15ijSjhbpK4QYFzJXh9LChE3eVhLFTXQzBYnRwDE1m5nGGV0IzkLLmV9N476w-OxOFy1Jkr3g-QM97azjU3Xerg8PpQkN6EU-9lPFqbkbZch5NZMe0cjptkr_6kmvJgzOZ_aqMIspC5mL94KwE_nXqHGRTlP8F7-B4MZJeciMQCNdCLQLnr9xggX6-j8Zl8Yyb-Zv9BR0ZRPRUVuIZQyBC4KXKvnODEFcXsXYINVrT7bEURDmEXpRJVV9MMFpltz2XwWuU1JFFVWhTfd68V2Lh2IOWta2Fd_k7rOg" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'

  
  curl -X POST https://your-project.vercel.app/v1/chat/completions \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer test-token" \
    -d '{
      "model": "gpt-4o",
      "messages": [{"role": "user", "content": "Hello!"}],
      "stream": false
    }'

    
    curl -X POST https://you2api-deploy.vercel.app//v1/chat/completions \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer  eyJhbGciOiJSUzI1NiIsImtpZCI6IlNLMmpJbnU3SWpjMkp1eFJad1psWHBZRUpQQkFvIiwidHlwIjoiSldUIn0.eyJhbXIiOlsiZW1haWwiXSwiYXV0aDBJZCI6Imdvb2dsZS1vYXV0aDJ8MTAxMTMwNDgwMTAzMzgzMzcxNDIzIiwiY3JlYXRlVGltZSI6MTcyMjMwNzcxMywiZHJuIjoiRFMiLCJlbWFpbCI6ImxpdWppZTE5OTYzMTBAZ21haWwuY29tIiwiZXhwIjoxNzU2NjUyMDI1LCJnaXZlbk5hbWUiOiJKaWUiLCJpYXQiOjE3NTU0NDI0MjUsImlzcyI6IlAyakludHRSTXVYcHlZWk1iVmNzYzRDOVowUlQiLCJsYXN0TmFtZSI6IkxpdSIsImxvZ2luSWRzIjpbImxpdWppZTE5OTYzMTBAZ21haWwuY29tIl0sIm5hbWUiOiIiLCJwaWN0dXJlIjoiIiwicmV4cCI6IjIwMjYtMDgtMTZUMTQ6NTM6NDVaIiwic3R5dGNoSWQiOiJ1c2VyLWxpdmUtYmE4OWQxNzctZDgxZi00YzYzLWJlNjMtMWYzNTUwZjI0ZjlkIiwic3ViIjoiVTJqd290TTZQN0o0cmZJczhvZzdGZjBVaDZIUiIsInN1YnNjcmlwdGlvblRpZXIiOiIiLCJ0ZW5hbnRDdXN0b21BdHRyaWJ1dGVzIjp7ImlzRW50ZXJwcmlzZSI6Int7dGVuYW50LmN1c3RvbUF0dHJpYnV0ZXMuaXNFbnRlcnByaXNlfX0iLCJuYW1lIjoie3t0ZW5hbnQubmFtZX19In0sInRlbmFudEludml0YXRpb24iOiIiLCJ0ZW5hbnRJbnZpdGVyIjoiIiwidXNlcklkIjoiVTJqd290TTZQN0o0cmZJczhvZzdGZjBVaDZIUiIsInZlcmlmaWVkRW1haWwiOnRydWV9.JN2yaXeRMP9IOh19grD6MVdTgj81Vlj2tIbmw4PomeQ5yJB_Y15ijSjhbpK4QYFzJXh9LChE3eVhLFTXQzBYnRwDE1m5nGGV0IzkLLmV9N476w-OxOFy1Jkr3g-QM97azjU3Xerg8PpQkN6EU-9lPFqbkbZch5NZMe0cjptkr_6kmvJgzOZ_aqMIspC5mL94KwE_nXqHGRTlP8F7-B4MZJeciMQCNdCLQLnr9xggX6-j8Zl8Yyb-Zv9BR0ZRPRUVuIZQyBC4KXKvnODEFcXsXYINVrT7bEURDmEXpRJVV9MMFpltz2XwWuU1JFFVWhTfd68V2Lh2IOWta2Fd_k7rOg" \
      -d '{"model": "gpt-4o", "messages": [{"role": "user", "content": "Hello!"}], "stream": false}'
      
      
      
      curl -X POST https://you2api-deploy-git-main-liujiers-projects.vercel.app/test \
        -H "Content-Type: application/json" \
        -d '{
          "messages": [
            {
              "role": "user", 
              "content": "Hello, how are you?"
            }
          ]
        }'

