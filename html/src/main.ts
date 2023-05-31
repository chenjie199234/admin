import { createApp } from 'vue'
import main from './main.vue'
// import { createVuestic } from "vuestic-ui"
// import "vuestic-ui/css"
// createApp(main).use(createVuestic()).mount('#app')
import {
	createVuesticEssential,
	VaSplit,
	VaModal,
	VaModalPlugin,
	VaCard,
	VaCardTitle,
	VaCardContent,
	VaHover,
	VaButton,
	VaSelect,
	VaInput,
	VaRadio,
	VaSwitch,
	VaImage,
	VaIcon,
	VaDivider,
	VaPagination,
	VaDropdown,
	VaDropdownContent,
	VaPopover,
	VaDropdownPlugin,
} from 'vuestic-ui'
import 'vuestic-ui/css'

createApp(main).use(createVuesticEssential({
		components: { VaSplit, VaModal, VaCard, VaCardTitle, VaCardContent, VaHover, VaButton, VaSelect, VaInput, VaRadio, VaSwitch, VaImage, VaIcon, VaDivider, VaPagination, VaDropdown, VaDropdownContent, VaPopover },
		plugins: { VaModalPlugin, VaDropdownPlugin },
		config: {},
})).mount('#app')
