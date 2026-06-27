import SiteHeader from './components/SiteHeader'
import HeroSection from './components/HeroSection'
import ProblemSection from './components/ProblemSection'
import IntegrationSection from './components/IntegrationSection'
import FeaturesSection from './components/FeaturesSection'
import SkillsSection from './components/SkillsSection'
import QuickStartSection from './components/QuickStartSection'
import SiteFooter from './components/SiteFooter'

function App() {
  return (
    <div style={{ minHeight: '100vh', background: '#0a0a0a' }}>
      <SiteHeader />
      <main>
        <HeroSection />
        <ProblemSection />
        <IntegrationSection />
        <FeaturesSection />
        <SkillsSection />
        <QuickStartSection />
      </main>
      <SiteFooter />
    </div>
  )
}

export default App
