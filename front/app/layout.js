import './globals.css'
import ContextProvider from './context'

export const metadata = {
  title: 'ByFood BookIT',
  description: 'An assignment in golang, NextJS, DBs and APIs',
}

export default function RootLayout({ children }) {

  return (
    <html lang="en">
      <body>
        <ContextProvider>
          <h1 style={{ textAlign: "center", fontFamily: "helvetica", fontSize: "3em", backgroundColor: "white", color: "#ff4D55" }}>ByFood BookIT</h1>
          {children}
        </ContextProvider>
      </body>
    </html>
  )
}
